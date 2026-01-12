package badgetemplates

// summaryTemplate is the SVG template for the Summary badge style (400x200)
const DefaultSummary = `<svg
  width="400"
  height="200"
  viewBox="0 0 400 200"
  xmlns="http://www.w3.org/2000/svg"
  role="img"
  aria-label="Open Source Contribution Summary">
  <defs>
    <style>
      text {
        font-family: system-ui, -apple-system, BlinkMacSystemFont,
                     "Segoe UI", Helvetica, Arial, sans-serif;
      }
      .bg {
        fill: {{.Colors.Background}};
      }
      .card {
        fill: {{.Colors.BackgroundAlt}};
      }
      .username {
        font-size: 18px;
        font-style: italic;
        font-weight: 700;
        fill: {{.Colors.Text}};
      }
      .subtitle {
        font-size: 12px;
        fill: {{.Colors.TextSecondary}};
      }
      .stat-value {
        font-size: 26px;
        font-style: italic;
        font-weight: 700;
        fill: {{.Colors.Text}};
      }
      .stat-label {
        font-size: 11px;
        fill: {{.Colors.TextSecondary}};
      }
    </style>
  </defs>
  <!-- Background -->
  <rect class="bg" width="400" height="200" rx="16"/>
  <!-- Header -->
  <text class="username" x="28" y="42">@{{.Stats.Username}}</text>
  <text class="subtitle" x="28" y="62">open source contributions</text>
  <!-- Stat Cards -->
  <rect class="card" x="22" y="91" width="108" height="70" rx="10"/>
  <rect class="card" x="146" y="91" width="108" height="70" rx="10"/>
  <rect class="card" x="270" y="91" width="108" height="70" rx="10"/>
  <!-- Stats -->
  <text class="stat-value" x="76" y="123" text-anchor="middle">{{.TotalProjects}}</text>
  <text class="stat-label" x="76" y="144" text-anchor="middle">PROJECTS</text>
  <text class="stat-value" x="200" y="123" text-anchor="middle">{{.TotalPRs}}</text>
  <text class="stat-label" x="200" y="144" text-anchor="middle">PRS MERGED</text>
  <text class="stat-value" x="324" y="123" text-anchor="middle">{{.TotalLines}}</text>
  <text class="stat-label" x="324" y="144" text-anchor="middle">LINES CHANGED</text>
</svg>
`

// compactTemplate is the SVG template for the Compact badge style (280x28) - Shields.io style
const DefaultCompact = `<svg
  width="280"
  height="32"
  viewBox="0 0 280 32"
  xmlns="http://www.w3.org/2000/svg"
  role="img"
  aria-label="OSS Contributions">

  <defs>
    <style>
      text {
        font-family: system-ui, -apple-system, BlinkMacSystemFont,
                     "Segoe UI", Helvetica, Arial, sans-serif;
        fill: {{.Colors.Text}};
        font-size: 12px;
        font-weight: 700;
      }

      .card {
        fill: {{.Colors.Background}};
        stroke: {{.Colors.Border}};
        stroke-width: 1;
        rx: 16;
      }
    </style>

    <linearGradient id="badgeGradient" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0%" stop-color="{{.Colors.Accent}}" stop-opacity="0.95"/>
      <stop offset="100%" stop-color="{{.Colors.Accent}}" stop-opacity="0.75"/>
    </linearGradient>
  </defs>

  <!-- Background -->
  <rect class="card" x="0.5" y="0.5" width="279" height="31"/>

  <!-- Text -->
  <text x="140" y="21" text-anchor="middle">
    OSS · {{.CompactText}}
  </text>
</svg>`

// detailedTemplate is the SVG template for the Detailed badge style (400x320)
const DefaultDetailed = `
<svg
  width="900"
  height="{{add 278 (mul (div (add (len .TopContributions) 2) 3) 120)}}"
  viewBox="0 0 900 {{add 278 (mul (div (add (len .TopContributions) 2) 3) 120)}}"
  fill="none"
  xmlns="http://www.w3.org/2000/svg"
  role="img"
  aria-label="GitHub Open Source Contribution Stats">

  <defs>
    <!-- Card gradient (solid fill with opacity) -->
    <linearGradient id="cardGradient" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0%" stop-color="{{.Colors.BackgroundAlt}}" stop-opacity="0.8"/>
      <stop offset="100%" stop-color="{{.Colors.BackgroundAlt}}" stop-opacity="0.6"/>
    </linearGradient>

    <!-- Glass overlay gradient -->
    <linearGradient id="glassOverlay" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0%" stop-color="{{.Colors.BackgroundAlt}}" stop-opacity="0.2"/>
      <stop offset="100%" stop-color="{{.Colors.BackgroundAlt}}" stop-opacity="0"/>
    </linearGradient>
  </defs>

  <!-- ========================= -->
  <!-- Background -->
  <!-- ========================= -->
  <rect
    width="900"
    height="{{add 278 (mul (div (add (len .TopContributions) 2) 3) 120)}}"
    rx="18"
    fill="{{.Colors.Background}}"
  />

  <!-- ========================= -->
  <!-- Header -->
  <!-- ========================= -->
  <text
    x="33"
    y="49"
    fill="{{.Colors.Text}}"
    font-family="Inter, system-ui, -apple-system, sans-serif"
    font-size="26"
    font-weight="bold"
    letter-spacing="0em">{{.Stats.Username}} · Open Source Contributions</text>

  <!-- ========================= -->
  <!-- Metrics Row -->
  <!-- ========================= -->

  <!-- Projects Card -->
  <g>
    <rect x="32" y="96" width="260" height="96" rx="14" fill="url(#cardGradient)" stroke="{{.Colors.Border}}"/>
    <rect x="32" y="96" width="260" height="96" rx="14" fill="url(#glassOverlay)"/>
    <text
      x="48"
      y="132"
      fill="{{.Colors.TextSecondary}}"
      font-family="Inter, system-ui, -apple-system, sans-serif"
      font-size="12"
      letter-spacing="0em">PROJECTS</text>
    <text
      x="48"
      y="167"
      fill="{{.Colors.Text}}"
      font-family="Inter, system-ui, -apple-system, sans-serif"
      font-size="28"
      font-weight="bold"
      letter-spacing="0em">{{.TotalProjects}}</text>
  </g>

  <!-- PRs Merged Card -->
  <g>
    <rect x="320" y="96" width="260" height="96" rx="14" fill="url(#cardGradient)" stroke="{{.Colors.Border}}"/>
    <rect x="320" y="96" width="260" height="96" rx="14" fill="url(#glassOverlay)"/>
    <text
      x="336"
      y="132"
      fill="{{.Colors.TextSecondary}}"
      font-family="Inter, system-ui, -apple-system, sans-serif"
      font-size="12"
      letter-spacing="0em">PRS MERGED</text>
    <text
      x="336"
      y="167"
      fill="{{.Colors.Text}}"
      font-family="Inter, system-ui, -apple-system, sans-serif"
      font-size="28"
      font-weight="bold"
      letter-spacing="0em">{{.TotalPRs}}</text>
  </g>

  <!-- Lines Changed Card -->
  <g>
    <rect x="608" y="96" width="260" height="96" rx="14" fill="url(#cardGradient)" stroke="{{.Colors.Border}}"/>
    <rect x="608" y="96" width="260" height="96" rx="14" fill="url(#glassOverlay)"/>
    <text
      x="624"
      y="132"
      fill="{{.Colors.TextSecondary}}"
      font-family="Inter, system-ui, -apple-system, sans-serif"
      font-size="12"
      letter-spacing="0em">LINES CHANGED</text>
    <text
      x="624"
      y="167"
      fill="{{.Colors.Text}}"
      font-family="Inter, system-ui, -apple-system, sans-serif"
      font-size="28"
      font-weight="bold"
      letter-spacing="0em">{{.TotalLines}}</text>
  </g>

  <!-- ========================= -->
  <!-- Top Contributions -->
  <!-- ========================= -->
  <text
    x="32"
    y="224"
    fill="{{.Colors.TextSecondary}}"
    font-family="Inter, system-ui, -apple-system, sans-serif"
    font-size="13"
    letter-spacing="0em">Top Contributions</text>

  {{range $i, $r := .TopContributions}}
  {{$col := mod $i 3}}
  {{$row := div $i 3}}
  {{$x := add 32 (mul $col 288)}}
  {{$y := add 240 (mul $row 120)}}
  <g>
    <!-- Card background -->
    <rect
      x="{{$x}}"
      y="{{$y}}"
      width="260"
      height="108"
      rx="14"
      fill="url(#cardGradient)"
      stroke="{{$.Colors.Border}}"
    />
    <!-- Glass overlay -->
    <rect
      x="{{$x}}"
      y="{{$y}}"
      width="260"
      height="108"
      rx="14"
      fill="url(#glassOverlay)"
    />
    <!-- Repo name -->
    <text
      x="{{add $x 16}}"
      y="{{add $y 32}}"
      fill="{{$.Colors.Text}}"
      font-family="Inter, system-ui, -apple-system, sans-serif"
      font-size="16"
      font-weight="bold"
      letter-spacing="0em">{{$r.RepoName}}</text>
    <!-- Owner -->
    <text
      x="{{add $x 16}}"
      y="{{add $y 50}}"
      fill="{{$.Colors.TextSecondary}}"
      font-family="Inter, system-ui, -apple-system, sans-serif"
      font-size="12"
      letter-spacing="0em">@{{$r.Owner}}</text>
    <!-- Stats -->
    <text
      x="{{add $x 16}}"
      y="{{add $y 90}}"
      fill="{{$.Colors.TextSecondary}}"
      font-family="Inter, system-ui, -apple-system, sans-serif"
      font-size="12"
      letter-spacing="0em">⭐ {{$r.Stars}} · {{$r.PRs}} PRs Merged</text>
  </g>
  {{end}}

</svg>
`
