# DKG 통합 보안 감사 리포트

**날짜**: 2026-03-07
**감사 범위**: Story CL (`client/x/dkg/`), story-kernel (`service/`, `store/`, `enclave/`, `server/`), DKG 컨트랙트 (`DKG.sol`, `SGXValidationHook.sol`)
**감사 방법**: 코드 정적 분석, 프로토콜 흐름 시뮬레이션, 암호학적 검증
**감사 대상**: DKG Setup, Resharing, Upgrade Resharing, Threshold Decryption 전체 워크플로우

---

## 목차

1. [감사 요약](#1-감사-요약)
2. [아키텍처 개요](#2-아키텍처-개요)
3. [Resharing 동작성 시뮬레이션](#3-resharing-동작성-시뮬레이션)
4. [발견 사항 — CRITICAL](#4-발견-사항--critical)
5. [발견 사항 — HIGH](#5-발견-사항--high)
6. [발견 사항 — MEDIUM](#6-발견-사항--medium)
7. [발견 사항 — LOW / INFO](#7-발견-사항--low--info)
8. [테스트 커버리지](#8-테스트-커버리지)
9. [우선순위별 권고사항](#9-우선순위별-권고사항)

---

## 1. 감사 요약

| 심각도 | 건수 |
|--------|------|
| **CRITICAL** | 2 |
| **HIGH** | 12 |
| **MEDIUM** | 14 |
| **LOW** | 12 |
| **INFO** | 5 |
| **합계** | **45** |

### 핵심 위험 요약

- **Resharing 동작성**: Validator set 변경 시 old-only 멤버의 Phase 전이 결함으로 dealing 실패 가능 (C-01)
- **TEE 보안 경계**: DKG state가 평문으로 저장되어 비밀 share 노출 위험 (C-02)
- **Upgrade Resharing**: 실패 시 복구 경로 없음, upgrade info 유실 (H-08)
- **인증 부재**: gRPC 서버에 TLS/인증 없음, MsgAddDkgVote에 authority 검증 없음

---

## 2. 아키텍처 개요

### 2.1 DKG 라이프사이클

```
Registration → Dealing → Finalization → Active → (resharing 시) Registration...
```

### 2.2 세 계층의 역할

| 계층 | 역할 | 주요 파일 |
|------|------|----------|
| **DKG Contract** | 등록/완료 이벤트, SGX attestation 검증 | `DKG.sol`, `SGXValidationHook.sol` |
| **Story CL** | 상태 전이, Vote Extension, 이벤트 핸들링 | `abci.go`, `dkg_svc_*.go`, `vote.go` |
| **story-kernel** | DKG 프로토콜 실행, 키 생성/관리, TDH2 | `dkg_service.go`, `dist_key_gen.go` |

### 2.3 Resharing 유형

| 유형 | 트리거 | IsResharing | IsUpgrade |
|------|--------|-------------|-----------|
| **일반 Resharing** | Active Period 만료 | true | false |
| **Upgrade Resharing** | Kernel Upgrade Activation | true | true |

---

## 3. Resharing 동작성 시뮬레이션

### 3.1 시나리오 A: 일반 Resharing (동일 Validator Set)

**가정**: Round 1 (Active) → Round 2 시작, Validator {A, B, C} 동일

#### Step 1: BeginBlocker에서 Stage 전이

```
abci.go: elapsed >= activeEnd
  → shouldTransitionStage returns (DKGStageRegistration, true)
  → 하지만 latestRound.Stage가 먼저 DKGStageRegistration으로 변경됨 (line 64)
  → 그 다음 InitiateDKGRound(ctx, false) 호출 (line 75)
```

**문제 발견**: `latestRound.Stage = nextStage` (line 64)로 **현재 active round의 stage를 Registration으로 변경**한 후 `setDKGNetwork`를 호출합니다. 그 다음 `InitiateDKGRound`에서 **새로운** round를 생성합니다. 이전 round의 stage가 Active에서 Registration으로 바뀌는 것은 의도치 않은 동작일 수 있습니다.

→ 코드 확인 결과: `DKGStageActive` → `DKGStageRegistration` 전이 시 `InitiateDKGRound`이 새 round를 생성하므로, 이전 round의 stage 변경은 무의미합니다. **기능적 영향 없음, 하지만 데이터 정합성 문제** (M-13으로 기록).

#### Step 2: shouldReshare 판정

```
dkg_initialization.go: shouldReshare(ctx)
  → getLatestActiveDKGNetwork(ctx) → Round 1 (Active)
  → activeNetwork != nil → return true
```

**정상**: IsResharing=true 설정됨.

#### Step 3: Registration

```
handleDKGRegistration:
  → isInCurRoundSet = slices.Contains(Round2.ActiveValSet, validatorAddr) → true (동일 set)
  → NewDKGSession(round=2, isResharing=true)
  → callTEEGenerateAndSealKey → GenerateAndSealKey (new key 생성)
  → callContractRegister → DKG.sol register() 호출
  → session.UpdatePhase(PhaseInitialized)
```

**정상**: 모든 validator가 현재 set에 있으므로 키 생성 정상 진행.

**이후 Phase 전이**: `PhaseInitialized`에서 `PhaseDealing`로의 전이는 어디서 발생하는가?

→ 코드 추적: `BeginDealing` (dkg_dealing.go:11) → `handleDKGDealing` 호출 → `session.Phase != types.PhaseDealing` 체크 (dkg_svc_dealing.go:57)

**문제**: `PhaseInitialized`에서 `PhaseDealing`로의 명시적 전이 코드가 없습니다. `handleDKGDealing`에서 `session.Phase != PhaseDealing`이면 스킵됩니다.

→ **추가 확인**: `ResumeDKGService` (dkg_svc.go:57-69)에서 `DKGStageDealing` 시 `session.UpdatePhase(types.PhaseDealing)`을 호출합니다. 그러나 이것은 **PhaseFailed 상태에서만** 호출됩니다 (line 37: `session.Phase != types.PhaseFailed`이면 return).

→ 실제로는 `BeginDealing`이 호출될 때 emitBeginDKGDealing 이벤트 발생 → 이 이벤트를 story-kernel이 받아서 Phase를 PhaseDealing로 변경하는 경로가 있어야 합니다.

→ **결론**: `handleDKGDealing`에서 `session.Phase`가 `PhaseInitialized`인 경우는 `PhaseDealing`이 아니므로 스킵됩니다. 이것은 **Phase 전이 경로가 불명확**한 문제입니다.

→ **재확인**: `dkg_svc_dealing.go`의 `handleDKGDealing`은 `Phase != PhaseDealing`이면 MarkFailed + return합니다. Phase가 `PhaseInitialized`이면 dealing이 실패합니다.

→ **Resolution**: Phase 전이는 어디서 일어나는가? 코드 전체를 추적한 결과:

1. `handleDKGRegistration` → `session.UpdatePhase(PhaseInitialized)` (완료)
2. **Phase transition to PhaseDealing은 명시적으로 없음**
3. `handleDKGDealing`에서 `Phase == PhaseDealing`을 기대함

→ 이것은 `ResumeDKGService`가 매 블록마다 호출되어 Phase를 업데이트해야 하는 설계이거나, registration 완료 후 Phase를 PhaseDealing으로 자동 전이해야 합니다.

→ **확인**: `ResumeDKGService`는 `Phase == PhaseFailed`일 때만 동작합니다. `PhaseInitialized`에서 `PhaseDealing`로의 전이가 **존재하지 않습니다**.

→ **이것은 C-01 (Old-only 멤버)뿐만 아니라, 모든 멤버에게 해당하는 문제일 수 있습니다.**

→ **최종 확인**: Phase 전이 로직을 다시 확인합니다. `handleDKGDealing` (dkg_svc_dealing.go:18-116)에서:

```go
session, err := k.stateManager.GetSession(dkgNetwork.Round)
if session.Phase != types.PhaseDealing {
    log.Warn(ctx, "Session not in dealing phase, skipping generate deals")
    k.stateManager.MarkFailed(ctx, session)
    return
}
```

→ 만약 이 코드가 실행 시점에 `PhaseInitialized` 상태이면 dealing이 실행되지 않습니다.

→ **그러나** 실제 동작에서는 DKG가 작동하고 있으므로, **Phase 전이가 다른 경로로 발생할 수 있습니다**. `stateManager`의 내부 로직이나, 이벤트 기반 Phase 업데이트가 있을 수 있습니다.

→ 코드를 재확인한 결과: `handleDKGRegistration`에서 `session.UpdatePhase(types.PhaseInitialized)` 후, **BeginDealing이 호출되기 전에** Phase를 PhaseDealing으로 전이하는 코드가 필요합니다.

→ **이 부분은 `BeginDealing` 자체에서** Phase 전이를 수행해야 합니다. 현재 `BeginDealing`은 Phase를 변경하지 않고 바로 `handleDKGDealing`을 goroutine으로 실행합니다.

→ **가설**: `PhaseInitialized` 이후 외부에서 Phase를 `PhaseDealing`로 바꾸는 별도의 이벤트 핸들러가 있을 수 있습니다. 하지만 코드에서 그런 핸들러는 발견되지 않았습니다.

→ **최종 판단**: **이것은 모든 resharing 시나리오에 영향을 미치는 CRITICAL 이슈입니다.** 단, 실제로 DKG가 작동한다면 이 Phase 전이가 다른 경로(이벤트 기반 등)로 처리되고 있을 수 있으므로, **실제 테스트로 확인 필요합니다.** C-01로 기록합니다.

#### Step 4: Dealing

```
shouldDeal(ctx, dkgNetwork):
  → inPrevSet = isInPrevActiveValSet(ctx) → Round 1의 ActiveValSet에 A,B,C 있음 → true
  → return true (resharing: prev set만 deal)

handleDKGDealing:
  → session.Phase == PhaseDealing 확인 (위의 Phase 전이 문제 해당)
  → dealerCC = session.CodeCommitment (일반 resharing이므로 OldCodeCommitment override 없음)
  → GenerateDeals → kyber GetResharingPrevDKG → DistKeyGenerator.Deals()
  → EnqueueDeals(resp.GetDeals())
```

**kernel 측**: `GenerateDeals`에서 `IsResharing=true`이므로:
```go
distKeyGen, err = s.GetResharingPrevDKG(cc, round, threshold, nextPubs, latest)
```
→ `GetResharingPrevDKG`는 이전 round의 DistKeyShare를 사용하여 deal 생성.

**정상**: (Phase 전이 문제가 해결된다는 전제 하에)

#### Step 5: ProcessDeals

```
handleDKGProcessDeals:
  → isInCurRoundSet 확인 → A,B,C 모두 current set → 처리
  → deal.RecipientIndex == session.Index인 deal만 필터링
  → kernel.ProcessDeals → GetResharingNextDKG → ProcessDeal 실행
  → EnqueueResponses
```

**kernel 측**: `ProcessDeals`에서 `IsResharing=true`이므로:
```go
distKeyGen, err = s.GetResharingNextDKG(cc, round, threshold, nextPubs)
```
→ Next committee DKG에서 deal 처리.

**정상**: Next committee가 deal을 받는 것은 Pedersen resharing에서 올바름.

#### Step 6: ProcessResponses

```
handleDKGProcessResponses:
  → shouldProcessResponses: inCurSet || inPrevSet → true
  → filteredResponses: resp.VssResponse.Index != session.Index인 것만
  → ccsToProcess: [session.CodeCommitment] (upgrade 아니므로 단일)
  → kernel.ProcessResponses
```

**kernel 측**: `ProcessResponses`에서 `IsResharing=true`이므로:
```go
prevDistKeyGen = GetResharingPrevDKG(...)
nextDistKeyGen = GetResharingNextDKG(...)
// 두 DKG 인스턴스에 모두 response 처리
```

**정상**: Prev + Next 모두에 response를 적용하는 것은 올바름.

#### Step 7: Finalization

```
handleDKGFinalization:
  → callTEEFinalizeDKG → kernel.FinalizeDKG
  → kernel: GetResharingNextDKG → DistKeyShare() → 새 share 생성
  → SealAndStoreDistKeyShare
  → callContractFinalizeDKG → DKG.sol finalize()
```

**정상**: Next committee의 DistKeyShare로 finalization.

#### Step 8: FinalizeDKGRound

```
dkg_finalization.go: FinalizeDKGRound
  → finalizedCount >= threshold 확인
  → settleRewardsForPreviousCommittee
  → setLatestActiveRound(ctx, latestRound) → Round 2가 active
```

**판정**: **동일 Validator Set 일반 Resharing은 Phase 전이 문제(C-01)를 제외하면 동작 가능.**

---

### 3.2 시나리오 B: Validator Set 변경 Resharing

**가정**: Round 1의 set {A, B, C}, Round 2의 set {B, C, D}
- A: old-only (빠짐)
- D: new-only (새로 참여)
- B, C: 양쪽 모두

#### Validator A (old-only):

**Registration**:
```
isInCurRoundSet = false (A는 Round 2의 ActiveValSet에 없음)
→ session 생성됨 (NOTE 주석: old members도 session 필요)
→ 키 생성 스킵 (line 75-79)
→ session.Phase = PhaseInitializing (UpdatePhase 미호출)
→ return (Phase 전이 없이 종료)
```

**Dealing**:
```
shouldDeal: inPrevSet = true (A는 Round 1의 set에 있음) → deal해야 함
handleDKGDealing: session.Phase == PhaseInitializing ≠ PhaseDealing → 스킵!
```

**→ A는 deal을 생성하지 못함 (C-01 확인)**

#### Validator D (new-only):

**Registration**:
```
isInCurRoundSet = true
→ 키 생성 (GenerateAndSealKey)
→ contract register
→ session.Phase = PhaseInitialized
```

**Dealing**:
```
shouldDeal: inPrevSet = false (D는 Round 1에 없음),
            collections.ErrNotFound가 아닌 경우 → return inPrevSet = false

→ shouldDeal = false → D는 deal 생성 안 함 (정상: new-only는 deal 생성 안 함)
```

**ProcessDeals**:
```
isInCurRoundSet = true → deal 수신/처리
→ Phase가 PhaseDealing이어야 함 (C-01 문제 동일)
```

**판정**: **C-01로 인해 old-only 멤버의 dealing 불가. Phase 전이 문제가 모든 멤버에게도 영향.**

---

### 3.3 시나리오 C: Upgrade Resharing

**가정**: Round 1 (old binary CC_old), Upgrade activation → Round 2 (new binary CC_new)

#### Step 1: Activation

```
abci.go:
  → hasPendingUpgradeActivation: upgradeInfo != nil, currentHeight >= activationHeight
  → DeleteKernelUpgradeInfo(ctx, upgradeInfo.UpgradeVersion) ← ⚠️ 여기서 삭제
  → InitiateDKGRound(ctx, true) ← isUpgrade=true
```

**문제 (H-08)**: Upgrade info가 activation 시점에 즉시 삭제됨. Round가 실패하면 복구 불가.

#### Step 2: Registration

```
handleDKGRegistration:
  → session.IsUpgrade = true
  → oldCC = getOldCodeCommitment(ctx) → Round 1 registration의 CC
  → session.OldCodeCommitment = oldCC

  For validators in current set:
    → getRegistrationKernelClient(ctx, isUpgrade=true, oldCC)
    → allCCs에서 oldCC와 다른 CC 찾음 → new binary client
    → GenerateAndSealKey via new binary
    → session.CodeCommitment = CC_new (new binary의 CC)
```

#### Step 3: Dealing

```
handleDKGDealing:
  → shouldDeal: inPrevSet = true (prev set의 멤버만)
  → dealerCC = session.OldCodeCommitment (IsUpgrade이므로 old CC 사용)
  → kernelRouter.GetClient(oldCC) → old binary client
  → old binary: GenerateDeals(IsResharing=true) → GetResharingPrevDKG
    → fromRound = latest.Round (Round 1)
    → old binary에 Round 1의 DistKeyShare 있음 → deal 생성
```

**정상**: (Phase 전이 문제 제외)

#### Step 4: ProcessDeals

```
handleDKGProcessDeals:
  → current set만 처리
  → session.CodeCommitment = CC_new → new binary로 라우팅
  → new binary: ProcessDeals(IsResharing=true) → GetResharingNextDKG
```

**정상**

#### Step 5: ProcessResponses

```
handleDKGProcessResponses:
  → ccsToProcess = [CC_new, CC_old] (dual-binary)
  → new binary: ProcessResponses → prev + next DKG
  → old binary: ProcessResponses → prev + next DKG

  Old binary의 ProcessResponses에서:
    → GetResharingPrevDKG: fromRound = latest.Round (Round 1)
      → old binary에 Round 1 share 있음 → OK
    → GetResharingNextDKG: round = Round 2
      → old binary에 Round 2의 키가 없을 수 있음 (new binary가 키 생성)
      → LoadSealedEd25519Key(CC_old, Round 2) 실패 → ❌
```

**문제 (H-04)**: Old binary에서 Round 2의 `GetResharingNextDKG` 실행 시 키가 없으면 실패. 에러가 발생해도 `continue`로 넘어가므로 old binary 쪽 DKG 상태 불일치.

#### Step 6: Finalization

```
handleDKGFinalization:
  → session.CodeCommitment = CC_new → new binary로 라우팅
  → new binary: FinalizeDKG(IsResharing=true) → GetResharingNextDKG → DistKeyShare()
```

**조건부 정상**: New binary가 충분한 deal + response를 받았으면 성공.

#### 실패 시

```
FinalizeDKGRound: finalizedCount < threshold
  → SkipToNextRound(ctx, currentRound)
  → InitiateDKGRound(ctx, false) ← isUpgrade=false!
  → upgrade info는 이미 삭제됨 → 복구 불가 (H-08)
```

**판정**: **Upgrade Resharing은 복수의 이슈 중첩으로 가장 높은 위험. 특히 실패 시 복구 경로 없음.**

---

### 3.4 시나리오 D: Kernel 재시작 후 Resharing 재개

```
story-kernel 재시작:
  → 모든 캐시 초기화 (DKGCache, ResharingCache, DistKeyShareCache 비어짐)
  → GenerateDeals 호출 시 GetResharingPrevDKG:
    → 캐시 miss → rebuildResharingPrevDKG
    → LoadDKGState(CC, fromRound) → 상태 파일 로드
    → rebuildResharingNextDKG(CC, fromRound):
      → LoadDKGState(CC, toRound) → stNext.FromRound = 0 (M-10: 직렬화 누락)
      → HasDKGState(CC, 0) → false
      → fallback to on-chain query (GetLatestActiveDKGNetwork)
    → replayMessages → 에러 무시 (L-09)
```

**판정**: **FromRound 누락으로 rebuild 경로가 on-chain fallback에 의존. 네트워크 지연/오류 시 실패.**

---

## 4. 발견 사항 — CRITICAL

### C-01: Phase 전이 결함 — PhaseInitialized에서 PhaseDealing으로의 전이 경로 없음

**위치**: `dkg_svc_registration.go:95`, `dkg_svc_dealing.go:57`

**문제**: `handleDKGRegistration`이 완료되면 session은 `PhaseInitialized` 상태입니다. 그러나 `handleDKGDealing`은 `Phase == PhaseDealing`을 요구합니다. `PhaseInitialized` → `PhaseDealing`으로의 명시적 전이 코드가 없습니다.

`ResumeDKGService`는 `PhaseFailed` 상태에서만 Phase를 복구합니다. 정상 흐름에서의 전이가 빠져있습니다.

**영향**: **모든 DKG round의 dealing이 실패할 수 있음**. 실제로 동작하고 있다면 다른 경로로 Phase가 전이되는 것이지만, 코드상으로는 확인되지 않음.

Old-only 멤버의 경우 추가로: 키 생성을 건너뛰므로 `PhaseInitializing`에서도 멈춤.

**권고**:
1. `BeginDealing` 또는 `handleDKGDealing` 시작부에서 `PhaseInitialized` → `PhaseDealing` 전이 로직 추가
2. Old-only 멤버: registration 시 `PhaseInitialized` → `PhaseDealing` 전이 추가

---

### C-02: story-kernel DKG State가 평문으로 저장됨 (비밀 Share 노출)

**위치**: `story-kernel/store/dkg_state.go:88-103`

**문제**: `saveState()`에서 `os.WriteFile(path, bz, 0o600)`로 DKG 상태를 저장합니다. 이 상태에는 `dkg.Deal` (비밀 share 포함)이 포함되어 있습니다. `key_store.go`는 올바르게 `enclave.SealToFile()`을 사용하지만, `dkg_state.go`는 평문 JSON을 사용합니다.

**영향**: TEE 호스트 관리자가 파일시스템 접근으로 비밀 share를 탈취 가능. threshold-1개의 share를 모으면 마스터 키 복원 가능.

**권고**: `dkg_state.go`도 `enclave.SealToFile()` 또는 `SealedLevelDB`를 사용하도록 변경.

---

## 5. 발견 사항 — HIGH

### H-01: gRPC 서버에 TLS/인증 없음

**위치**: `story-kernel/server/server.go:27`

**문제**: `grpc.NewServer()`로 기본 서버 생성. TLS, mTLS, 토큰 인증 등 없음. TEE 내부 통신이라 해도 같은 호스트의 다른 프로세스가 DKG 명령을 주입할 수 있음.

**권고**: mTLS 또는 Unix domain socket + 파일 권한으로 접근 제한.

### H-02: Vote Extension 중복 제거 미구현

**위치**: `vote.go:84` — `// TODO: discard duplicated votes`

**문제**: 동일한 deal/response가 여러 validator의 VE에 포함될 경우 모두 그대로 집계됨. 악의적 validator가 동일 deal을 반복 전파 가능.

**권고**: `(sessionID, dealerIndex, recipientIndex)` 튜플 기준 중복 제거 구현.

### H-03: MsgAddDkgVote Authority 검증 누락

**위치**: `msg_server.go:16-47`

**문제**: `AddVote` 핸들러에서 `msg.Authority` 검증이 없음. `latestRound`과 `Stage == DKGStageDealing` 확인만 수행. 모듈 주소가 아닌 임의의 주소에서 트랜잭션 제출 가능.

**권고**: `msg.Authority == authtypes.NewModuleAddress(types.ModuleName).String()` 검증 추가.

### H-04: Upgrade Resharing에서 Old Binary의 ProcessResponses/ProcessJustifications 실패

**위치**: `dkg_svc_dealing.go:233-300`, `dkg_svc_dealing.go:432-471`

**문제**: Upgrade resharing 시 responses/justifications를 old + new binary 모두에 전달하지만, old binary에서 Round N+1의 키가 없으면 `GetResharingNextDKG` 실패. 에러가 `continue`로 무시됨.

**영향**: Old binary의 DKG 상태 불일치. Finalization 시 `GetResharingPrevDKG`에서 DistKeyShare 실패 가능.

**권고**: Old binary 처리 실패 시 명확한 에러 로깅 및 재시도 전략 필요.

### H-05: 전역 가변 큐 사용 (Deal/Response/Justification)

**위치**: `keeper.go:24-29`, `queue.go`

**문제**: Deal, response, justification 큐가 `var` 패키지 레벨 전역 변수로 관리됨. Keeper 인스턴스와 무관하며, 노드 재시작 시 데이터 유실.

**권고**: Keeper 구조체 내부로 이동. 영속성이 필요하면 KV store 사용.

### H-06: VerifyVoteExtension에서 Semantic 검증 부재

**위치**: `vote.go:38-46`

**문제**: Proto unmarshal 성공 여부만 확인. Deal/response/justification의 세션 ID, 인덱스 범위, 라운드 번호 등 semantic 검증 없음.

**권고**: 최소한 round 번호, 인덱스 범위 검증 추가.

### H-07: DKG.sol finalize()에 호출자 검증 불충분

**위치**: `DKG.sol:210-241`

**문제**: `finalize()`는 `chargesFee` + `whenNotPaused` modifier만 있음. 수수료를 낼 수 있는 누구나 호출 가능. Validator가 아닌 주소에서 잘못된 global public key를 제출할 수 있음.

**분석**: `validatorAddr` 파라미터로 validator를 지정하지만, `msg.sender == validatorAddr` 확인이 없음. 제3자가 다른 validator 명의로 finalize를 호출할 수 있음.

**권고**: `require(msg.sender == validatorAddr)` 또는 validator set 교차 검증 추가.

### H-08: Upgrade Round 실패 시 Upgrade Info 유실 — 복구 불가

**위치**: `abci.go:54`, `dkg_round.go:55`

**문제**:
1. `hasPendingUpgradeActivation` 시 `DeleteKernelUpgradeInfo` 호출 (abci.go:54)
2. Round 실패 시 `SkipToNextRound` → `InitiateDKGRound(ctx, false)` — isUpgrade=false (dkg_round.go:55)
3. Upgrade info가 삭제되었으므로 재시도 불가

**영향**: Upgrade resharing이 한 번 실패하면 영구적으로 실패. 새로운 `scheduleUpgrade()` 트랜잭션을 다시 제출해야 함.

**권고**: Upgrade info를 round 성공 시에만 삭제하거나, 실패 시 `isUpgrade=true`를 유지.

### H-09: DKG.sol register()에 Validator Set 검증 없음

**위치**: `DKG.sol:171-199`

**문제**: SGX attestation만 검증하고, 호출자가 실제 bonded validator인지 확인하지 않음. 비-validator가 DKG에 등록하여 프로토콜 방해 가능.

**분석**: Story CL의 `Registered()` handler (dkg_handler.go:46)에서 `slices.Contains(latest.ActiveValSet, validator)` 확인이 있음. 하지만 컨트랙트 레벨에서는 누구나 등록 가능하며 이벤트가 발생함.

**영향**: 가스 낭비 공격, 이벤트 스팸.

**권고**: 컨트랙트에서도 validator set 검증 추가, 또는 precompile을 통한 validator 확인.

### H-10: SGXValidationHook에 whenNotPaused 누락

**위치**: `SGXValidationHook.sol:79`

**문제**: DKG 컨트랙트가 pause되어도 `validateReport()`는 `msg.sender == DKG` 체크만 수행. DKG가 pause 상태이면 register/finalize가 호출되지 않으므로 실질적 영향은 제한적이나, 직접 호출에 대한 방어 없음.

### H-11: gRPC Reflection이 프로덕션에 활성화됨

**위치**: `story-kernel/server/server.go:39`

**문제**: `reflection.Register(svr)` — TODO 주석이 있으나 프로덕션 코드에 포함. 공격자가 서비스 인터페이스를 탐색 가능.

### H-12: 비밀 키 메모리 제로화 없음

**위치**: `story-kernel/service/dkg_service.go` 전반

**문제**: 비밀 share, 개인 키가 사용 후 메모리에서 명시적으로 제로화되지 않음.

---

## 6. 발견 사항 — MEDIUM

### M-01: DKGState.FromRound 필드가 Disk 직렬화에서 누락

**위치**: `story-kernel/store/dkg_state.go:35-41`

**문제**: `dkgStateDisk` 구조체에 `FromRound` 필드 없음. 저장 후 로드하면 `FromRound = 0`. Kernel 재시작 시 resharing DKG rebuild 실패 가능.

**권고**: `dkgStateDisk`에 `from_round` JSON 필드 추가.

### M-02: Resharing 시 Old Committee 참여율 사전 검증 없음

**위치**: `dkg_dealing.go:36-37`

**문제**: Threshold는 new set 기준으로만 계산. Old committee에서 threshold 미만 참여 시 finalization에서야 실패 감지.

### M-03: dkgSvcRunning Global Lock이 Stage 전환을 차단할 수 있음

**위치**: 각 handler의 `dkgSvcRunning.CompareAndSwap(false, true)`

**문제**: 하나의 handler가 실행 중이면 다른 stage의 handler 스킵. Stage 전환 시 한 블록 지연 가능.

### M-04: ActivePeriod 만료와 Upgrade Activation 동시 발생 시 우선순위

**위치**: `abci.go:41-86`

**문제**: Upgrade activation이 항상 우선. Active period 만료에 의한 일반 resharing 스킵. 문서화 필요.

### M-05: SkipToNextRound에서 IsUpgrade 플래그 전파 안 됨

**위치**: `dkg_round.go:55`

**문제**: `SkipToNextRound` → `InitiateDKGRound(ctx, false)`. 이전 round가 upgrade이었어도 다음 round는 일반 round.

### M-06: 등록 시 공개키 충돌 미검사 (Rogue-Key 가능성)

**위치**: `dkg_handler.go:23-92`

**문제**: `Registered()` handler에서 동일 공개키 중복 등록 체크 없음.

### M-07: getDKGRegistrationsByRound에서 전체 Walk

**위치**: `dkg_registration.go:73-86`

**문제**: 모든 registration을 walk하면서 prefix 매칭. O(N) 스캔. Validator 수 증가 시 성능 저하.

### M-08: DKGNetwork 키 순서 비결정적 가능성

**위치**: `dkg_initialization.go:22-34`

**문제**: `GetActiveValidators()`에서 validator 순서가 `GetAllValidators()` 반환 순서에 의존. 노드 간 순서가 다를 수 있음.

**분석**: Cosmos SDK의 `GetAllValidators()`는 store iteration 순서 (deterministic)를 따르므로 실제로는 문제없을 수 있음. 하지만 명시적 정렬 없음.

### M-09: TDH2 PID Base Convention 불일치 위험

**위치**: `story-kernel/service/dkg_generate_deals.go:42-55`

**문제**: PID가 1-based (DKG registration index)로 설정되지만, kyber DKG는 0-based index 사용. 변환 오류 시 잘못된 share 매핑.

### M-10: Magic Prefix 0x04 0x3f 하드코딩

**위치**: `story-kernel/service/dkg_service.go:747`

**문제**: `buildTDH2PublicKey`에서 `append([]byte{0x04, 0x3f}, dkgPubKey...)`로 prefix 하드코딩. Edwards25519 → TDH2(secp256k1) 변환 시 매직 바이트 사용.

### M-11: Report Data 오프셋 하드코딩 (SGXValidationHook)

**위치**: `SGXValidationHook.sol:143-163`

**문제**: SGX Quote 구조의 오프셋(64, 368)이 하드코딩. Quote 형식 변경 시 깨짐.

### M-12: 수수료 소각 — 인출 불가

**위치**: `DKG.sol:34-38`

**문제**: `payable(address(0x0)).transfer(msg.value)` — 수수료가 0x0로 전송되어 영구 소각. 의도적일 수 있으나, 잘못된 수수료 설정 시 복구 불가.

### M-13: Active Stage에서 Registration으로 전이 시 이전 Round Stage 오염

**위치**: `abci.go:64-65`

**문제**: `latestRound.Stage = nextStage` 후 `setDKGNetwork`으로 저장. Active round의 stage가 Registration으로 변경됨. 새 round가 별도로 생성되므로 기능 영향 없지만 데이터 정합성 문제.

### M-14: Scalar 바이트 역순 변환 수동 구현

**위치**: `story-kernel/service/dkg_service.go:743`

**문제**: kyber ↔ cb-mpc 간 little-endian/big-endian 변환이 `reverseBytes()`로 수동 구현. 경계 조건 테스트 부족.

---

## 7. 발견 사항 — LOW / INFO

### LOW

| ID | 설명 | 위치 |
|----|------|------|
| L-01 | ProcessDeal 에러 시 silent drop (continue) | story-kernel/service/dkg_service.go:304 |
| L-02 | 세션 ID 기반 replay 방지 미구현 | story-kernel |
| L-03 | HKDF salt가 nil | dkg_service.go:787 |
| L-04 | ProcessResponse 에러 시 silent drop | dkg_service.go:406 |
| L-05 | 캐시 eviction 미구현 (unbounded growth) | store/dkg_cache.go |
| L-06 | replayMessages에서 에러 무시 | dist_key_gen.go:454-467 |
| L-07 | Justification 처리 시 partial failure 핸들링 없음 | dkg_justification.go |
| L-08 | 로그에 민감 정보 포함 가능 | 다수 파일 |
| L-09 | TEE quote 만료 시간 미검증 | enclave/quote.go |
| L-10 | 최대 참여자 수 제한 없음 | dkg_registration.go |
| L-11 | Finalize 이벤트 중복 처리 가능 | dkg_handler.go |
| L-12 | Key resharing 시 이전 키 폐기 로직 없음 | dkg_service.go |

### INFO

| ID | 설명 | 위치 |
|----|------|------|
| I-01 | dkgStartBlock=10 하드코딩 | abci.go:13 |
| I-02 | dkgAsyncTimeout=1분 — 충분한지 검토 필요 | dkg_svc.go:18 |
| I-03 | DKG 테스트 커버리지 극히 낮음 | contracts/test/ |
| I-04 | scheduleUpgrade/cancelUpgrade에 whenNotPaused 없음 (의도적일 수 있음) | DKG.sol |
| I-05 | CachePID에서 err가 nil인데 Wrap에 전달 | dkg_generate_deals.go:53 |

---

## 8. 테스트 커버리지

| 저장소 | 현재 커버리지 (추정) | 권장 |
|--------|---------------------|------|
| Story CL (`client/x/dkg/`) | ~15-20% | 최소 70% |
| DKG Contracts | ~5% (testDKG_Initialize 1건) | 최소 80% |
| story-kernel | ~10-15% | 최소 70% |

### 필수 테스트 목록

#### Resharing 관련
- [ ] Phase 전이 (PhaseInitializing → PhaseInitialized → PhaseDealing → PhaseFinalized → PhaseCompleted)
- [ ] `shouldDeal` + `shouldProcessResponses` resharing 조합
- [ ] Old-only 멤버의 dealing 참여
- [ ] Validator set 변경 시나리오 (add/remove validator)
- [ ] Upgrade resharing 전체 flow mock
- [ ] Upgrade round 실패 → 재시도 경로

#### Kernel 관련
- [ ] `DKGState` FromRound 직렬화 round-trip
- [ ] Kernel 재시작 후 DKG rebuild (resharing context)
- [ ] `GetResharingPrevDKG` / `GetResharingNextDKG` 캐시 miss rebuild
- [ ] Dual-binary routing (old + new CC)

#### Contract 관련
- [ ] register() — validator set 교차 검증
- [ ] finalize() — double finalization 방지
- [ ] scheduleUpgrade() → activationHeight → resharing 트리거
- [ ] Fee mechanism edge cases

---

## 9. 우선순위별 권고사항

### P0 — 즉시 조치

| 순서 | 이슈 | 조치 |
|------|------|------|
| 1 | C-01 | Phase 전이 로직 추가 (PhaseInitialized → PhaseDealing, old-only 멤버 포함) |
| 2 | C-02 | `dkg_state.go`에 SGX sealing 적용 |
| 3 | H-01 | gRPC 서버에 mTLS 또는 Unix socket 인증 추가 |
| 4 | H-03 | MsgAddDkgVote authority 검증 추가 |
| 5 | H-08 | Upgrade info를 round 성공 시에만 삭제 |

### P1 — 단기 조치

| 순서 | 이슈 | 조치 |
|------|------|------|
| 6 | M-01 | `dkgStateDisk`에 `FromRound` 필드 추가 |
| 7 | H-02 | Vote Extension 중복 제거 구현 |
| 8 | H-05 | 전역 큐를 keeper 내부로 이동 |
| 9 | H-06 | VerifyVoteExtension에 semantic 검증 추가 |
| 10 | H-07 | finalize()에 msg.sender == validatorAddr 검증 |
| 11 | H-04 | Dual-binary ProcessResponses 오류 처리 강화 |
| 12 | H-11 | 프로덕션에서 gRPC reflection 제거 |

### P2 — 중기 조치

| 순서 | 이슈 | 조치 |
|------|------|------|
| 13 | M-02~M-14 | Medium 이슈 순차 수정 |
| 14 | - | 전체 테스트 커버리지 70% 이상으로 향상 |
| 15 | - | DKG 컨트랙트 외부 감사 권장 |

### P3 — 장기 조치

| 순서 | 이슈 | 조치 |
|------|------|------|
| 16 | - | DKG 프로토콜 formal verification 검토 |
| 17 | - | TDH2 DLEQ proof 온체인 검증 구현 |
| 18 | - | Key rotation 및 proactive secret sharing 도입 |
| 19 | - | Resharing 전체 flow E2E 테스트 프레임워크 구축 |

---

## 종합 판정

### 동작성

| 시나리오 | 판정 | 핵심 위험 |
|---------|------|----------|
| 초기 DKG Setup | **C-01 확인 필요** | Phase 전이 |
| 동일 Set 일반 Resharing | **C-01 확인 필요** | Phase 전이 |
| Set 변경 Resharing | **위험** | C-01 + old-only 멤버 |
| Upgrade Resharing | **가장 위험** | C-01 + H-04 + H-08 중첩 |
| Kernel 재시작 후 재개 | **위험** | M-01 (FromRound 누락) |

### 보안

전체적으로 **DKG 프로토콜의 암호학적 설계는 올바르게 구현**되었습니다 (Pedersen DKG, Feldman VSS, kyber v4 사용). 그러나 **구현 레벨의 보안 경계**(인증, 권한 검증, 비밀 보호, Phase 관리)에서 다수의 취약점이 존재합니다.

가장 시급한 것은 **C-01 (Phase 전이)**과 **C-02 (평문 상태 저장)**입니다. C-01은 resharing 동작 자체를 무력화할 수 있고, C-02는 TEE의 보안 모델을 무력화합니다.
