# TUI Dashboard Layout Spec v2

## 핵심 원칙
1. **한 화면에 모든 정보가 보여야 함** — 스크롤 없이 24줄 터미널에서도 동작
2. **정렬이 완벽해야 함** — 모든 패널의 좌우 끝이 일치, 빈 공간 없음
3. **각 요소가 뭔지 즉시 구분 가능** — 라벨, 구분선, 색상으로 명확히
4. **반응형** — 창 크기 바꿔도 깨지지 않음

## Terminal Size
- **Minimum**: 80 cols × 24 rows
- **Target**: 120+ cols × 30+ rows
- **Responsive**: width < 100이면 세로 스택

---

## Layout (width >= 100)

```
 [1] server-name      [2] rpi5      (2 available · Tab to switch)
╭──────────────────────────────────╮╭────────────────────────────────────────────╮
│  ⚡ server-name                  ││  Docker Containers                         │
│                                  ││                                            │
│  CPU  ██████░░░░░░░░░░░░  17.1%  ││  NAME            STATE     IMAGE    STATUS │
│  Mem  ██████████████░░░░  46.5%  ││  homebridge      running   hb/hb    Up 12d │
│  /    █░░░░░░░░░░░░░░░░░   3.0%  ││  portainer       running   port~    Up 12d │
│                                  ││  pihole          exited    pihole   Ex 3d  │
│  ── History (2min) ────────────  ││  prometheus      running   prom~    Up 5d  │
│  CPU ▁▂▃▄▅▃▂▁▃▂▃▄▅▃▂▁▃▂        ││  ... and 3 more                            │
│  Mem ▅▆▆▅▆▆▅▆▅▆▆▅▆▆▅▆▅▆        ││                                            │
│                                  ││  Top Processes (CPU)                       │
│  Uptime:  42d 7h                 ││                                            │
│  OS:      darwin/arm64           ││  PID    CPU%   MEM%  NAME                  │
│  Cores:   10                     ││  462    8.3%   0.2%  WallpaperAerials~     │
│  Memory:  7.9 / 16.0 GB         ││  169    5.1%   0.4%  WindowServer          │
│                                  ││  561    2.7%   0.2%  VTDecoderXPCSer~     │
│                                  ││  392    0.9%   0.8%  iTerm2               │
│                                  ││  4679   0.6%   3.9%  openclaw-gateway     │
╰──────────────────────────────────╯╰────────────────────────────────────────────╯
╭───────────────────────────────────────────────────────────────────────────────╮
│  Alerts: CPU: 17%  Mem: 47%  Disk /: 3%                                      │
│  Tab/Shift+Tab switch server  │  q quit  │  ⟳ 2s                             │
╰───────────────────────────────────────────────────────────────────────────────╯
```

---

## 정렬 규칙 (반드시 지킬 것)

### 패널 폭
- **좌측 패널(System)**: `width * 2/5`
- **우측 패널(Docker+Processes)**: `width - leftWidth - 2` (border 고려)
- **Footer**: `leftWidth + rightWidth` (상단 패널 합산과 정확히 동일)
- 패널 사이 갭: **0** (border가 바로 붙음)

### 프로그레스 바 정렬
- 라벨 고정폭 5자: `"  CPU"`, `"  Mem"`, `"  /  "`
- 바와 퍼센트 사이 공백 1칸
- 퍼센트 고정폭 6자: `%5.1f%%`
- **barWidth = panelInnerWidth - labelWidth(5) - space(2) - percentWidth(7)**

### Sparkline 정렬
- History 섹션 내에서 라벨 4자: `"  CPU"`, `"  Mem"`
- sparkline 폭 = barWidth와 동일
- **lipgloss.Width()로 강제 고정폭 렌더링**

---

## 패널 상세

### Tab Bar
- 패널 밖, 최상단
- Active: bold + 배경색 #62
- Inactive: dim #241
- 우측: `(N available · Tab to switch)`

### Left Panel: System

**섹션 1 — Metrics (바)**
```
  CPU  ██████░░░░░░░░░░░░  17.1%
  Mem  ██████████████░░░░  46.5%
  /    █░░░░░░░░░░░░░░░░░   3.0%
```
- 바 사이 빈 줄 없음
- 색상: <70% green, 70-90% yellow, >=90% red

**섹션 2 — History (sparkline)**
```
  ── History (2min) ────────────
  CPU ▁▂▃▄▅▃▂▁▃▂▃▄▅▃▂▁▃▂
  Mem ▅▆▆▅▆▆▅▆▅▆▆▅▆▆▅▆▅▆
```
- 위에 구분선 (dimStyle)
- sparkline 색상: 마지막 값 기준 green/yellow/red
- 데이터 없으면: `▁` 반복 (dim)
- max 60 points, 왼쪽부터 자라남
- **바와 sparkline 사이에 빈 줄 1개로 구분**

**섹션 3 — Info**
```
  Uptime:  42d 7h
  OS:      darwin/arm64
  Cores:   10
  Memory:  7.9 / 16.0 GB
```
- sparkline과 info 사이 빈 줄 1개

### Right Panel: Docker + Processes (하나의 패널)

**Docker Containers (상단)**
```
  Docker Containers

  NAME            STATE     IMAGE    STATUS
  homebridge      running   hb/hb    Up 12d
  ...
  ... and 3 more
```
- 최대 8개, 초과 시 "... and N more"
- running=green, exited=red, other=yellow
- Docker 없으면: `Docker not installed` (dim)
- Docker 안 되면: `Docker unavailable` (yellow)

**Top Processes (하단, 같은 패널 안)**
```
  Top Processes (CPU)

  PID    CPU%   MEM%  NAME
  462    8.3%   0.2%  WallpaperAerials~
  ...
```
- Docker와 같은 패널 안에서 빈 줄 2개로 구분
- 최대 5개
- 프로세스 이름 truncate with ~

### Footer
```
  Alerts: CPU: 17%  Mem: 47%  Disk /: 3%
  Tab/Shift+Tab switch server  │  q quit  │  ⟳ 2s
```
- 전체 폭 = 상단 좌+우 패널 합산
- Alerts 색상: ok=green, warning=yellow, critical=red

---

## 색상 팔레트
| 용도 | 색상코드 | 예시 |
|------|---------|------|
| Panel border | #62 (blue) | 패널 테두리 |
| Active tab bg | #62 | 선택된 탭 |
| Title | #230 (bold) | ⚡ server-name |
| OK/Green | #82 | running, 정상 수치 |
| Warning/Yellow | #226 | 70-90% |
| Critical/Red | #196 | >=90%, exited |
| Dim | #241 | 비활성, 구분선 |
| Header | #62 (bold) | 테이블 헤더 |

## 바/Sparkline 문자
- Progress bar: `█` (filled) + `░` (empty)
- Sparkline: `▁▂▃▄▅▆▇█` (8 levels, 0-100%)

---

## 데이터 규칙
- **CPU %**: 시스템 전체, 0-100% 캡 (load average 기반이라도 100 초과 시 cap)
- **Process CPU%**: 캡 없음 (멀티코어에서 100% 초과 정상)
- **Refresh**: 2초
- **History**: 최대 60 points (~2분)
- **Docker 타임아웃**: 2초 (초과 시 캐시 반환)

---

## 반응형 (width < 100)
세로 스택:
```
[탭바]
[System 패널 - 전체폭]
[Docker 패널 - 전체폭]
[Processes 패널 - 전체폭]
[Footer - 전체폭]
```

---

## 구현 체크리스트
- [ ] 모든 패널 폭 계산을 하나의 함수에서 관리
- [ ] footer 폭 = left + right 패널 합산 (border 포함)
- [ ] sparkline을 History 섹션으로 분리 (구분선 + 라벨)
- [ ] Docker + Processes를 하나의 패널로 합침
- [ ] barWidth 계산: 패널 내부폭 - 라벨 - 퍼센트 - 여백
- [ ] lipgloss.Width()로 모든 가변폭 요소 강제 고정
- [ ] 테스트: 80, 100, 120, 138, 160 cols에서 렌더링 확인
