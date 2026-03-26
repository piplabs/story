//go:build integration

package dkg

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockCallRecord 记录一次 RPC 调用。
type MockCallRecord struct {
	Method string
	Round  uint32
	Err    error
}

// MockKernelServer 是可注入行为的 story-kernel gRPC mock server。
// 默认行为：所有 RPC 返回 Unimplemented，除非通过 On* 函数覆盖。
// 用于集成测试中模拟 TEE 故障注入、签名伪造等场景。
type MockKernelServer struct {
	dkgtypes.UnimplementedKernelServiceServer

	mu    sync.Mutex
	Calls []MockCallRecord

	// 行为注入：若非 nil 则覆盖默认行为
	OnGetCodeCommitment func(ctx context.Context, req *dkgtypes.GetCodeCommitmentRequest) (*dkgtypes.GetCodeCommitmentResponse, error)
	OnGenerateAndSealKey func(ctx context.Context, req *dkgtypes.GenerateAndSealKeyRequest) (*dkgtypes.GenerateAndSealKeyResponse, error)
	OnGenerateDeals      func(ctx context.Context, req *dkgtypes.GenerateDealsRequest) (*dkgtypes.GenerateDealsResponse, error)
	OnProcessDeals       func(ctx context.Context, req *dkgtypes.ProcessDealsRequest) (*dkgtypes.ProcessDealsResponse, error)
	OnProcessResponses   func(ctx context.Context, req *dkgtypes.ProcessResponsesRequest) (*dkgtypes.ProcessResponsesResponse, error)
	OnProcessJustification func(ctx context.Context, req *dkgtypes.ProcessJustificationRequest) (*dkgtypes.ProcessJustificationResponse, error)
	OnFinalizeDKG        func(ctx context.Context, req *dkgtypes.FinalizeDKGRequest) (*dkgtypes.FinalizeDKGResponse, error)
	OnPartialDecryptTDH2 func(ctx context.Context, req *dkgtypes.PartialDecryptTDH2Request) (*dkgtypes.PartialDecryptTDH2Response, error)

	// CodeCommitment 用于 GetCodeCommitment 的默认返回值
	CodeCommitment []byte
}

func (m *MockKernelServer) record(method string, round uint32, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = append(m.Calls, MockCallRecord{Method: method, Round: round, Err: err})
}

// GetCallCount 返回指定方法的调用次数。
func (m *MockKernelServer) GetCallCount(method string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	count := 0
	for _, c := range m.Calls {
		if c.Method == method {
			count++
		}
	}
	return count
}

// ResetCalls 清除所有调用记录。
func (m *MockKernelServer) ResetCalls() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = nil
}

func (m *MockKernelServer) GetCodeCommitment(ctx context.Context, req *dkgtypes.GetCodeCommitmentRequest) (*dkgtypes.GetCodeCommitmentResponse, error) {
	if m.OnGetCodeCommitment != nil {
		resp, err := m.OnGetCodeCommitment(ctx, req)
		m.record("GetCodeCommitment", 0, err)
		return resp, err
	}
	cc := m.CodeCommitment
	if cc == nil {
		cc = []byte("mock-code-commitment-for-testing")
	}
	m.record("GetCodeCommitment", 0, nil)
	return &dkgtypes.GetCodeCommitmentResponse{CodeCommitment: cc}, nil
}

func (m *MockKernelServer) GenerateAndSealKey(ctx context.Context, req *dkgtypes.GenerateAndSealKeyRequest) (*dkgtypes.GenerateAndSealKeyResponse, error) {
	if m.OnGenerateAndSealKey != nil {
		resp, err := m.OnGenerateAndSealKey(ctx, req)
		m.record("GenerateAndSealKey", req.Round, err)
		return resp, err
	}
	m.record("GenerateAndSealKey", req.Round, nil)

	dkgPubKey := make([]byte, 32)
	dkgPubKey[0] = 0xAA
	commPubKey := make([]byte, 64)
	commPubKey[0] = 0xBB

	// Build fake SGX quote with correct code/data commitments at the offsets
	// that SGXValidationHook.validateReport expects:
	//   quote[112:144] = MRENCLAVE (code commitment)
	//   quote[368:400] = report_data first 32 bytes (data commitment)
	// Requires a fake DCAP contract (always returns true) to bypass signature verification.
	fakeQuote := buildFakeSGXQuote(m.CodeCommitment, req.Address, req.Round, 1, make([]byte, 32), dkgPubKey, commPubKey)

	return &dkgtypes.GenerateAndSealKeyResponse{
		CodeCommitment:   m.CodeCommitment,
		Round:            req.Round,
		StartBlockHeight: 1,
		StartBlockHash:   make([]byte, 32),
		DkgPubKey:        dkgPubKey,
		CommPubKey:       commPubKey,
		EnclaveReport:    fakeQuote,
	}, nil
}

// buildFakeSGXQuote constructs a minimal 432-byte fake SGX quote with correct
// code commitment (MRENCLAVE) and data commitment at the offsets expected by
// SGXValidationHook.validateReport. Use with a fake DCAP contract that skips
// signature chain verification.
func buildFakeSGXQuote(codeCommitment []byte, validatorAddr string, round uint32, startBlockHeight int64, startBlockHash, dkgPubKey, commPubKey []byte) []byte {
	quote := make([]byte, 432)

	// code commitment at offset 112 (48-byte header + 64 bytes into report body)
	if len(codeCommitment) >= 32 {
		copy(quote[112:144], codeCommitment[:32])
	}

	// data commitment at offset 368 (48-byte header + 320 bytes into report body)
	// Must match DKG.sol: keccak256(validatorAddr || round || startBlockHeight || startBlockHash || dkgPubKey || enclaveCommKey)
	dataCommitment := computeDataCommitment(validatorAddr, round, uint64(startBlockHeight), startBlockHash, dkgPubKey, commPubKey)
	copy(quote[368:400], dataCommitment)

	return quote
}

// computeDataCommitment mirrors DKG.sol and story-kernel's calculateReportData:
// keccak256(abi.encodePacked(validatorAddr(20), round(4), startBlockHeight(8), startBlockHash(32), dkgPubKey, enclaveCommKey))
func computeDataCommitment(validatorAddr string, round uint32, startBlockHeight uint64, startBlockHash, dkgPubKey, commPubKey []byte) []byte {
	addr := common.HexToAddress(validatorAddr)

	roundBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(roundBytes, round)

	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, startBlockHeight)

	var buf []byte
	buf = append(buf, addr.Bytes()...)   // 20 bytes
	buf = append(buf, roundBytes...)     // 4 bytes
	buf = append(buf, heightBytes...)    // 8 bytes
	buf = append(buf, startBlockHash...) // 32 bytes
	buf = append(buf, dkgPubKey...)      // variable
	buf = append(buf, commPubKey...)     // variable

	return crypto.Keccak256(buf)
}

func (m *MockKernelServer) GenerateDeals(ctx context.Context, req *dkgtypes.GenerateDealsRequest) (*dkgtypes.GenerateDealsResponse, error) {
	if m.OnGenerateDeals != nil {
		resp, err := m.OnGenerateDeals(ctx, req)
		m.record("GenerateDeals", req.Round, err)
		return resp, err
	}
	m.record("GenerateDeals", req.Round, nil)
	return &dkgtypes.GenerateDealsResponse{
		CodeCommitment: m.CodeCommitment,
		Round:          req.Round,
	}, nil
}

func (m *MockKernelServer) ProcessDeals(ctx context.Context, req *dkgtypes.ProcessDealsRequest) (*dkgtypes.ProcessDealsResponse, error) {
	if m.OnProcessDeals != nil {
		resp, err := m.OnProcessDeals(ctx, req)
		m.record("ProcessDeals", req.Round, err)
		return resp, err
	}
	m.record("ProcessDeals", req.Round, nil)
	return &dkgtypes.ProcessDealsResponse{
		CodeCommitment: m.CodeCommitment,
		Round:          req.Round,
	}, nil
}

func (m *MockKernelServer) ProcessResponses(ctx context.Context, req *dkgtypes.ProcessResponsesRequest) (*dkgtypes.ProcessResponsesResponse, error) {
	if m.OnProcessResponses != nil {
		resp, err := m.OnProcessResponses(ctx, req)
		m.record("ProcessResponses", req.Round, err)
		return resp, err
	}
	m.record("ProcessResponses", req.Round, nil)
	return &dkgtypes.ProcessResponsesResponse{}, nil
}

func (m *MockKernelServer) ProcessJustification(ctx context.Context, req *dkgtypes.ProcessJustificationRequest) (*dkgtypes.ProcessJustificationResponse, error) {
	if m.OnProcessJustification != nil {
		resp, err := m.OnProcessJustification(ctx, req)
		m.record("ProcessJustification", req.Round, err)
		return resp, err
	}
	m.record("ProcessJustification", req.Round, nil)
	return &dkgtypes.ProcessJustificationResponse{}, nil
}

func (m *MockKernelServer) FinalizeDKG(ctx context.Context, req *dkgtypes.FinalizeDKGRequest) (*dkgtypes.FinalizeDKGResponse, error) {
	if m.OnFinalizeDKG != nil {
		resp, err := m.OnFinalizeDKG(ctx, req)
		m.record("FinalizeDKG", req.Round, err)
		return resp, err
	}
	m.record("FinalizeDKG", req.Round, nil)
	return &dkgtypes.FinalizeDKGResponse{
		CodeCommitment:   m.CodeCommitment,
		Round:            req.Round,
		ParticipantsRoot: make([]byte, 32),
		GlobalPubKey:     []byte("mock-global-pub-key"),
		PubKeyShare:      []byte("mock-pub-key-share"),
		Signature:        []byte("mock-signature"),
	}, nil
}

func (m *MockKernelServer) PartialDecryptTDH2(ctx context.Context, req *dkgtypes.PartialDecryptTDH2Request) (*dkgtypes.PartialDecryptTDH2Response, error) {
	if m.OnPartialDecryptTDH2 != nil {
		resp, err := m.OnPartialDecryptTDH2(ctx, req)
		m.record("PartialDecryptTDH2", req.Round, err)
		return resp, err
	}
	m.record("PartialDecryptTDH2", req.Round, nil)
	return &dkgtypes.PartialDecryptTDH2Response{
		EncryptedPartialDecryption: []byte("mock-encrypted-partial"),
		EphemeralPubKey:            []byte("mock-ephemeral-pub-key"),
		PubShare:                   []byte("mock-pub-share"),
		Signature:                  []byte("mock-signature"),
	}, nil
}

// ──────────────────────────────────────────────────────────────────────────────
// MockKernelInstance 管理一个运行中的 mock kernel gRPC server 实例。
// ──────────────────────────────────────────────────────────────────────────────

// MockKernelInstance 包装 gRPC server + listener，可启动和关闭。
type MockKernelInstance struct {
	Server     *MockKernelServer
	GRPCServer *grpc.Server
	Listener   net.Listener
	Addr       string
}

// StartMockKernel 在随机端口启动一个 mock kernel gRPC server。
func StartMockKernel(codeCommitment []byte) (*MockKernelInstance, error) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("listen: %w", err)
	}

	mockSrv := &MockKernelServer{
		CodeCommitment: codeCommitment,
	}
	grpcSrv := grpc.NewServer()
	dkgtypes.RegisterKernelServiceServer(grpcSrv, mockSrv)

	go func() {
		_ = grpcSrv.Serve(lis)
	}()

	return &MockKernelInstance{
		Server:     mockSrv,
		GRPCServer: grpcSrv,
		Listener:   lis,
		Addr:       lis.Addr().String(),
	}, nil
}

// Stop 停止 mock kernel server。
func (inst *MockKernelInstance) Stop() {
	inst.GRPCServer.GracefulStop()
}

// ──────────────────────────────────────────────────────────────────────────────
// 预定义行为注入工厂函数
// ──────────────────────────────────────────────────────────────────────────────

// WithForgedSignature 返回 PartialDecryptTDH2 handler 产出伪造签名的 partial。
func WithForgedSignature() func(context.Context, *dkgtypes.PartialDecryptTDH2Request) (*dkgtypes.PartialDecryptTDH2Response, error) {
	return func(_ context.Context, req *dkgtypes.PartialDecryptTDH2Request) (*dkgtypes.PartialDecryptTDH2Response, error) {
		return &dkgtypes.PartialDecryptTDH2Response{
			EncryptedPartialDecryption: []byte("forged-encrypted-partial"),
			EphemeralPubKey:            []byte("forged-ephemeral-pub-key"),
			PubShare:                   []byte("forged-pub-share"),
			Signature:                  []byte("FORGED-INVALID-SIGNATURE-12345"), // 伪造签名
		}, nil
	}
}

// WithWrongCommKey 返回 PartialDecryptTDH2 handler 用错误的 commPubKey 签名。
func WithWrongCommKey() func(context.Context, *dkgtypes.PartialDecryptTDH2Request) (*dkgtypes.PartialDecryptTDH2Response, error) {
	return func(_ context.Context, req *dkgtypes.PartialDecryptTDH2Request) (*dkgtypes.PartialDecryptTDH2Response, error) {
		return &dkgtypes.PartialDecryptTDH2Response{
			EncryptedPartialDecryption: []byte("wrong-key-encrypted-partial"),
			EphemeralPubKey:            []byte("wrong-ephemeral-pub-key"),
			PubShare:                   []byte("wrong-pub-share"),
			Signature:                  []byte("wrong-key-signature"), // 用不同 key 签名
		}, nil
	}
}

// WithWrongPubShare 返回 PartialDecryptTDH2 handler 产出不匹配的 pubShare。
func WithWrongPubShare() func(context.Context, *dkgtypes.PartialDecryptTDH2Request) (*dkgtypes.PartialDecryptTDH2Response, error) {
	return func(_ context.Context, req *dkgtypes.PartialDecryptTDH2Request) (*dkgtypes.PartialDecryptTDH2Response, error) {
		return &dkgtypes.PartialDecryptTDH2Response{
			EncryptedPartialDecryption: []byte("valid-encrypted-partial"),
			EphemeralPubKey:            []byte("valid-ephemeral-pub-key"),
			PubShare:                   []byte("MISMATCHED-PUB-SHARE-WRONG"), // 与 registration.pubKeyShare 不匹配
			Signature:                  []byte("valid-signature-but-wrong-share"),
		}, nil
	}
}

// WithDecryptError 返回 PartialDecryptTDH2 handler 直接返回错误（TEE 解密失败）。
func WithDecryptError(msg string) func(context.Context, *dkgtypes.PartialDecryptTDH2Request) (*dkgtypes.PartialDecryptTDH2Response, error) {
	return func(_ context.Context, _ *dkgtypes.PartialDecryptTDH2Request) (*dkgtypes.PartialDecryptTDH2Response, error) {
		return nil, status.Errorf(codes.Internal, "TEE decrypt failed: %s", msg)
	}
}

// WithOversizedDeals 返回 GenerateDeals handler 产出超大 deals（>256KB 或 >80 个）。
func WithOversizedDeals(count int) func(context.Context, *dkgtypes.GenerateDealsRequest) (*dkgtypes.GenerateDealsResponse, error) {
	return func(_ context.Context, req *dkgtypes.GenerateDealsRequest) (*dkgtypes.GenerateDealsResponse, error) {
		deals := make([]dkgtypes.Deal, count)
		for i := range count {
			deals[i] = dkgtypes.Deal{
				Index:          uint32(i),
				RecipientIndex: uint32(i % 3),
				Deal: dkgtypes.EncryptedDeal{
					DhKey:     make([]byte, 1024),
					Signature: make([]byte, 64),
					Nonce:     make([]byte, 12),
					Cipher:    make([]byte, 4096),
				},
				Signature: make([]byte, 64),
			}
		}
		return &dkgtypes.GenerateDealsResponse{
			CodeCommitment: req.CodeCommitment,
			Round:          req.Round,
			Deals:          deals,
		}, nil
	}
}

// WithGarbageDeals 返回 GenerateDeals handler 产出 malformed proto 数据。
func WithGarbageDeals() func(context.Context, *dkgtypes.GenerateDealsRequest) (*dkgtypes.GenerateDealsResponse, error) {
	return func(_ context.Context, req *dkgtypes.GenerateDealsRequest) (*dkgtypes.GenerateDealsResponse, error) {
		// 返回带损坏数据的 deals
		return &dkgtypes.GenerateDealsResponse{
			CodeCommitment: req.CodeCommitment,
			Round:          req.Round,
			Deals: []dkgtypes.Deal{{
				Index:          0,
				RecipientIndex: 0,
				Deal: dkgtypes.EncryptedDeal{
					DhKey:     []byte("GARBAGE-MALFORMED-DATA"),
					Signature: []byte("INVALID"),
					Nonce:     []byte("BAD"),
					Cipher:    []byte("CORRUPTED-CIPHER-DATA-NOT-VALID-AES-GCM"),
				},
				Signature: []byte("GARBAGE-SIGNATURE-NOT-VALID-SCHNORR"),
			}},
		}, nil
	}
}

// WithInvalidVSSDeal 返回 GenerateDeals handler 产出有效签名但无效 VSS share（触发 justification）。
func WithInvalidVSSDeal() func(context.Context, *dkgtypes.GenerateDealsRequest) (*dkgtypes.GenerateDealsResponse, error) {
	return func(_ context.Context, req *dkgtypes.GenerateDealsRequest) (*dkgtypes.GenerateDealsResponse, error) {
		// 签名结构正确但 VSS share 无效 → recipient 会 complaint → justification
		return &dkgtypes.GenerateDealsResponse{
			CodeCommitment: req.CodeCommitment,
			Round:          req.Round,
			Deals: []dkgtypes.Deal{{
				Index:          0,
				RecipientIndex: 1,
				Deal: dkgtypes.EncryptedDeal{
					DhKey:     make([]byte, 32),
					Signature: make([]byte, 64),
					Nonce:     make([]byte, 12),
					Cipher:    []byte("encrypted-but-invalid-vss-share"),
				},
				Signature: make([]byte, 64),
			}},
		}, nil
	}
}

// WithRoundReplay 返回各 RPC handler 重放指定 round 的数据（用于 cross-round replay 测试）。
// capturedDeals/capturedResponses 应在前一轮运行时捕获。
func WithRoundReplay(capturedRound uint32, capturedDeals []dkgtypes.Deal) func(context.Context, *dkgtypes.GenerateDealsRequest) (*dkgtypes.GenerateDealsResponse, error) {
	return func(_ context.Context, req *dkgtypes.GenerateDealsRequest) (*dkgtypes.GenerateDealsResponse, error) {
		// 无论当前 round 是多少，都返回旧 round 的 deals
		return &dkgtypes.GenerateDealsResponse{
			CodeCommitment: req.CodeCommitment,
			Round:          capturedRound, // 故意返回旧 round
			Deals:          capturedDeals,
		}, nil
	}
}

// WithFinalizeDKGError 返回 FinalizeDKG handler 直接返回错误（TEE finalize 失败）。
func WithFinalizeDKGError(msg string) func(context.Context, *dkgtypes.FinalizeDKGRequest) (*dkgtypes.FinalizeDKGResponse, error) {
	return func(_ context.Context, _ *dkgtypes.FinalizeDKGRequest) (*dkgtypes.FinalizeDKGResponse, error) {
		return nil, status.Errorf(codes.Internal, "TEE finalize failed: %s", msg)
	}
}

// WithGenerateAndSealKeyError 返回 GenerateAndSealKey handler 返回错误（session 创建失败）。
func WithGenerateAndSealKeyError(msg string) func(context.Context, *dkgtypes.GenerateAndSealKeyRequest) (*dkgtypes.GenerateAndSealKeyResponse, error) {
	return func(_ context.Context, _ *dkgtypes.GenerateAndSealKeyRequest) (*dkgtypes.GenerateAndSealKeyResponse, error) {
		return nil, status.Errorf(codes.Internal, "session creation failed: %s", msg)
	}
}
