# DKG Audit - 최종 결과 보고서

## 코드 수정 — 14개 항목 완료

| ID | 심각도 | 설명 | 상태 |
|----|--------|------|------|
| C-01 | Critical | `handleDKGDealing` Phase 전환: PhaseInitialized 체크 → PhaseDealing 업데이트 | ✅ |
| H-02 | High | deal/response 중복 제거 (`deduplicateDeals`, `deduplicateResponses` 추가) | ✅ |
| H-03 | High | `msg_server.go`, `proposal_server.go`에 authority check 추가 | ✅ |
| H-06 | High | Vote Extension semantic validation 추가 (크기/개수 제한) | ✅ |
| H-08 | High | upgrade info 삭제를 성공적 resharing 완료 시에만 수행 | ✅ |
| H-11 | High | gRPC reflection을 DebugMode 플래그 뒤에 게이팅 | ✅ |
| H-12 | High | 개인키 메모리 제로화 (`zeroBytes`, `zeroPrivateKey` + defer) | ✅ |
| M-01 | Medium | `dkgStateDisk`에 `FromRound` 필드 추가 (직렬화 갭 수정) | ✅ |
| M-05 | Medium | `SkipToNextRound`에서 isUpgrade 플래그 전파 | ✅ |
| M-09 | Medium | PID 인덱싱 컨벤션 주석 추가 (1-based, kyber 0-based 관계) | ✅ |
| M-10 | Medium | TDH2 magic prefix를 명명 상수로 정의 + 주석 | ✅ |
| M-13 | Medium | `DKGStageEnded` 추가, 이전 라운드 Active 유지 → 새 라운드 complete 시 Ended | ✅ |
| M-14 | Medium | `reverseBytes` 엣지 케이스 테스트 8개 추가 | ✅ |
| I-05 | Info | `CachePID`: `errors.Wrap(nil, ...)` → `errors.New(...)` | ✅ |
| L-08 | Low | justification 전체 로깅 → 인덱스만 로깅 (민감정보 제거) | ✅ |

## 조사 완료 — 의도된 설계이므로 수정 불필요

| ID | 설명 | 결론 | 근거 |
|----|------|------|------|
| C-02 | DKG state sealing | **불필요** | Deal은 ECDH 암호화, Response는 비밀 미포함, Justification 평문은 프로토콜 의도. VE 공개 전파 |
| H-04 | dual-binary ProcessResponses | **문제 없음** | 각 binary가 자신의 역할만 처리, 나머지 graceful skip |
| H-05 | queue 영속성 전략 | **현재 적절** | ResumeDKGService 재시도, VE 재전파, stage window 충분 |
| L-02 | session replay 방지 | **현재 충분** | kyber SessionID 바인딩(라운드별 새 키), stage gating, 중복 제거 |
| L-03 | HKDF nil salt | **수용 가능** | IKM이 ECDH 고엔트로피, ephemeral 키 매번 생성, info로 도메인 분리 |
| L-07 | justification partial failure | **의도된 동작** | 부분 성공 설계, threshold가 누락 보상, gRPC 에러로 재시도 가능 |
| M-02 | old member 미참여 slashing | **threshold 처리** | t+1 참여면 성공, 자동 slashing은 정직 validator 오처벌 위험 |
| M-03 | stage 전환 edge case | **대부분 복구** | MarkFailed + ResumeDKGService, 성공 후 skip은 다음 블록 재시작 |
| M-06 | 라운드별 키 유일성 | **보장됨** | 파일 경로 `{round}/{cc}/key_ed25519`로 분리 |
| M-08 | ValidatorsKey determinism | **결정론적** | IAVL KVStore 사전순 반복, ActiveValSet on-chain 합의 |
| M-11 | SGX quote offset | **안전** | MRENCLAVE/REPORTDATA 오프셋은 SGX 아키텍처 수준 고정, UUPS 업그레이드 가능 |
| I-02 | dkgAsyncTimeout=1min | **충분** | 최악 ~40초, 실패 시 Resume 재시도. 2분 상향 고려 가능 |

## 낮은 우선순위 — 별도 PR 권장

| ID | 설명 | 이유 |
|----|------|------|
| L-12 | 이전 라운드 sealed key 정리 | SGX 봉인 보호 중이라 긴급하지 않으나, 최소 노출 원칙 따라 cleanup 추가 권장 |
| L-05 | 캐시 eviction (keep last N rounds) | 라운드 키 기반 unbounded 캐시, 장기 운영 시 메모리 증가 가능 |

## 문서화 완료

| ID | 설명 | 상태 |
|----|------|------|
| H-01 | `client/x/dkg/docs/grpc_tls.md` (mTLS 도입 계획) | ✅ |
| M-04 | DKG module README에 upgrade resharing 우선순위 note | ✅ |

## 명시적으로 무시 (의도된 설계 또는 별도 작업)

H-05(global vars), H-07, H-09, H-10, I-01, I-03, I-04, L-01, L-04, L-06, L-09, L-10, L-11, M-07, M-12

## 테스트 커버리지

| 패키지 | 커버리지 |
|--------|---------|
| story/client/x/dkg/keeper | 32.4% |
| story/client/x/dkg/types | 1.2% |
| story-kernel/store | 42.2% |
| story-kernel/service | 빌드 불가 (cb-mpc CGO 의존성, CI 환경 필요) |
