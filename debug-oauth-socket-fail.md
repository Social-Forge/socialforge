# Debug Session: oauth-socket-fail
- **Status**: [OPEN]
- **Issue**: OAuth Google flow belum tervalidasi end-to-end, dan frontend kadang menerima `fetch failed` dengan `UND_ERR_SOCKET` saat request ke backend `127.0.0.1:8080`.
- **Debug Server**: http://127.0.0.1:7777/event
- **Log File**: .dbg/trae-debug-log-oauth-socket-fail.ndjson

## Reproduction Steps
1. Jalankan frontend SvelteKit dan backend API lokal.
2. Buka halaman signin/signup.
3. Klik tombol OAuth Google.
4. Ikuti redirect hingga callback kembali.
5. Amati apakah login sukses atau frontend menampilkan error / backend menutup socket.

## Hypotheses & Verification
| ID | Hypothesis | Likelihood | Effort | Evidence |
|----|------------|------------|--------|----------|
| A | Backend OAuth init/callback crash sehingga socket ditutup sebelum memberi response | High | Medium | Rejected: backend health normal via `192.168.100.98:8080`, belum ada indikasi crash handler |
| B | Callback URL/provider config mismatch, jadi flow balik ke alamat yang salah atau state session tidak cocok | High | Low | Confirmed: Google menolak `redirect_uri=http://localhost:5000/auth/google/callback` dengan `Error 400: redirect_uri_mismatch` |
| C | Frontend route OAuth tidak meneruskan cookie/session backend yang dibutuhkan goth saat callback | High | Medium | Inconclusive: callback belum tercapai karena gagal lebih awal di API base URL |
| D | Backend API tidak aktif / restart loop / panic saat endpoint `/auth/oauth/*` dipanggil | Medium | Low | Confirmed and mitigated: konflik loopback `127.0.0.1:8080` dihindari dengan memindahkan backend ke `8081`; request post-fix sudah tembus ke backend |
| E | Header/redirect handling di proxy fetch server-side SvelteKit memicu request gagal sebelum response terbaca | Medium | Medium | Rejected as primary cause: socket failure terjadi sebelum response backend app terbaca, dan target socket menunjuk ke `127.0.0.1:8080` |

## Log Evidence
- Instrumentation aktif di frontend OAuth start/callback, frontend API request wrapper, dan backend handler OAuth redirect/callback.
- Backend handler instrumentation sudah lolos kompilasi (`go test ./internal/handlers/...`).
- Pre-fix NDJSON:
  - `frontend/src/lib/server/api.ts:createApiRequest` mencatat `baseUrl=http://localhost:8080/api` lalu gagal dengan `UND_ERR_SOCKET` ke `remoteAddress=127.0.0.1 remotePort=8080`.
  - `frontend/src/routes/auth/[provider]/+server.ts` mencatat route OAuth start memang terpanggil sebelum gagal.
- Post-fix NDJSON:
  - `frontend/src/routes/auth/[provider]/+server.ts` mencatat OAuth start sukses memanggil backend dan menerima redirect Google dengan `redirect_uri=http://localhost:5000/auth/google/callback`.
  - Browser test dari `/signin` berhenti di halaman Google error `Access blocked: This app's request is invalid`.
  - Detail error Google menampilkan `Error 400: redirect_uri_mismatch` untuk `redirect_uri=http://localhost:5000/auth/google/callback`.
- Probe manual:
  - `curl http://localhost:8080/health` => `Empty reply from server`
  - `curl http://192.168.100.98:8080/health` => `200 OK`
  - `http://localhost:8081/health` => `200 OK`
- Port binding:
  - `127.0.0.1:8080` juga dipakai `com.docker.backend`
  - backend Go `main.exe` listen di `0.0.0.0:8080`

## Verification Conclusion
Perbandingan bukti:
- Pre-fix: OAuth start gagal di server-side fetch karena `localhost:8080` bentrok di loopback dan menghasilkan `UND_ERR_SOCKET`.
- Post-fix: OAuth start sudah berhasil sampai redirect ke Google melalui `localhost:8081`, jadi socket issue lokal tidak lagi menjadi blocker utama.
- Blocker aktif saat ini: konfigurasi Google OAuth belum mengizinkan callback `http://localhost:5000/auth/google/callback`, sehingga flow berhenti sebelum user bisa login atau consent.
