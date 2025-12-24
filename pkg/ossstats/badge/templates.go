package badge

// summaryTemplate is the SVG template for the Summary badge style (400x200)
const summaryTemplate = `<svg
  width="400"
  height="200"
  viewBox="0 0 400 200"
  xmlns="http://www.w3.org/2000/svg">

  <defs>
    <style>
      :root {
        --bg: {{.Colors.Background}};
        --text: {{.Colors.Text}};
        --text-muted: {{.Colors.TextSecondary}};
        --border: {{.Colors.Border}};
        --accent: {{.Colors.Accent}};
        --card-bg: rgba(255,255,255,0.03);
      }

      text {
        font-family: system-ui, -apple-system, BlinkMacSystemFont,
                     "Segoe UI", Helvetica, Arial, sans-serif;
      }

      .card {
        fill: var(--card-bg);
        stroke: var(--border);
        stroke-width: 1;
        rx: 12;
      }

      .glass {
        fill: url(#glass);
        rx: 12;
      }

      .accent {
        fill: var(--accent);
      }

      .title {
        font-size: 15px;
        font-weight: 600;
        fill: var(--text);
      }

      .label {
        font-size: 11px;
        fill: var(--text-muted);
        text-transform: uppercase;
        letter-spacing: 0.08em;
      }

      .value {
        font-size: 22px;
        font-weight: 700;
        fill: var(--text);
      }
    </style>

    <linearGradient id="glass" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0%" stop-color="rgba(255,255,255,0.06)" />
      <stop offset="100%" stop-color="rgba(255,255,255,0)" />
    </linearGradient>
  </defs>

  <!-- Card -->
  <rect x="0.5" y="0.5" width="399" height="199" class="card"/>
  <rect x="0.5" y="0.5" width="399" height="199" class="glass"/>
  <rect x="12" y="0.5" width="376" height="3" rx="2" class="accent"/>

  <!-- Header -->
  <text x="20" y="32" class="title">Open Source Stats</text>
  <text x="20" y="52" class="label">@{{.Stats.Username}}</text>

  <!-- Metrics -->
  <g transform="translate(20, 90)">
    <text class="value" x="0" y="0">{{.TotalProjects}}</text>
    <text class="label" x="0" y="18">Projects</text>
  </g>

  <g transform="translate(210, 90)">
    <text class="value" x="0" y="0">{{.TotalPRs}}</text>
    <text class="label" x="0" y="18">PRs Merged</text>
  </g>

  <g transform="translate(20, 145)">
    <text class="value" x="0" y="0">{{.TotalLines}}</text>
    <text class="label" x="0" y="18">Lines Changed</text>
  </g>
</svg>`

// compactTemplate is the SVG template for the Compact badge style (280x28) - Shields.io style
const compactTemplate = `<svg
  width="280"
  height="28"
  viewBox="0 0 280 28"
  xmlns="http://www.w3.org/2000/svg">

  <defs>
    <style>
      :root {
        --bg: {{.Colors.Background}};
        --text: {{.Colors.Text}};
        --border: {{.Colors.Border}};
        --accent: {{.Colors.Accent}};
        --card-bg: rgba(255,255,255,0.05);
      }

      text {
        font-family: system-ui, -apple-system, BlinkMacSystemFont,
                     "Segoe UI", Helvetica, Arial, sans-serif;
        font-size: 11px;
        font-weight: 600;
        fill: var(--text);
      }

      .card {
        fill: var(--card-bg);
        stroke: var(--border);
        stroke-width: 1;
        rx: 6;
      }

      .accent {
        fill: var(--accent);
      }
    </style>
  </defs>

  <!-- Card -->
  <rect x="0.5" y="0.5" width="279" height="27" class="card"/>
  <rect x="8" y="0.5" width="264" height="2" rx="1" class="accent"/>

  <!-- Text -->
  <text x="140" y="18" text-anchor="middle">
    {{.CompactText}}
  </text>
</svg>`

// minimalTemplate is the SVG template for the Minimal badge style (120x28)
const minimalTemplate = `<svg
  width="120"
  height="28"
  viewBox="0 0 120 28"
  xmlns="http://www.w3.org/2000/svg">

  <defs>
    <style>
      :root {
        --bg: {{.Colors.Background}};
        --text: {{.Colors.Text}};
        --border: {{.Colors.Border}};
        --accent: {{.Colors.Accent}};
        --card-bg: rgba(255,255,255,0.08);
      }

      text {
        font-family: system-ui, -apple-system, BlinkMacSystemFont,
                     "Segoe UI", Helvetica, Arial, sans-serif;
        font-size: 11px;
        font-weight: 600;
        fill: var(--text);
      }

      .card {
        fill: var(--card-bg);
        stroke: var(--border);
        stroke-width: 1;
        rx: 6;
      }

      .accent {
        fill: var(--accent);
      }
    </style>
  </defs>

  <!-- Card -->
  <rect x="0.5" y="0.5" width="119" height="27" class="card"/>
  <rect x="8" y="0.5" width="104" height="2" rx="1" class="accent"/>

  <!-- Text -->
  <text x="60" y="18" text-anchor="middle">
    {{.MinimalText}}
  </text>
</svg>`

// detailedTemplate is the SVG template for the Detailed badge style (400x320)
const detailedTemplate = `
<svg
  width="900"
  height="440"
  viewBox="0 0 900 440"
  xmlns="http://www.w3.org/2000/svg"
  role="img"
  aria-label="GitHub Open Source Contribution Stats">

  <!-- ========================= -->
  <!-- Theme + Styling -->
  <!-- ========================= -->
  <style>
    :root {
      --bg: {{.Colors.Background}};
      --text: {{.Colors.Text}};
      --text-muted: {{.Colors.TextSecondary}};
      --border: {{.Colors.Border}};
      --accent: {{.Colors.Accent}};
      --card-bg: rgba(255,255,255,0.03);
    }

    text {
      font-family: system-ui, -apple-system, BlinkMacSystemFont,
                   "Segoe UI", Helvetica, Arial, sans-serif;
    }

    .title {
      font-size: 24px;
      font-weight: 700;
      fill: var(--text);
      letter-spacing: -0.3px;
    }

    .subtitle {
      font-size: 13px;
      fill: var(--text-muted);
    }

    .card {
      fill: var(--card-bg);
      stroke: var(--border);
      stroke-width: 1;
      rx: 14;
    }

    .card-glass {
      fill: url(#glass);
      rx: 14;
    }

    .accent-strip {
      fill: var(--accent);
    }

    .metric-label {
      font-size: 11px;
      text-transform: uppercase;
      letter-spacing: 0.08em;
      fill: var(--text-muted);
    }

    .metric-value {
      font-size: 24px;
      font-weight: 700;
      fill: var(--text);
    }

    .repo-name {
      font-size: 15px;
      font-weight: 600;
      fill: var(--text);
    }

    .repo-meta {
      font-size: 12px;
      fill: var(--text-muted);
    }

    .fade-in {
      opacity: 0;
      animation: fadeUp 0.6s ease-out forwards;
    }

    @keyframes fadeUp {
      from { opacity: 0; transform: translateY(6px); }
      to   { opacity: 1; transform: translateY(0); }
    }
  </style>

  <!-- ========================= -->
  <!-- Definitions -->
  <!-- ========================= -->
  <defs>
    <linearGradient id="glass" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0%" stop-color="rgba(255,255,255,0.06)" />
      <stop offset="100%" stop-color="rgba(255,255,255,0)" />
    </linearGradient>
  </defs>

  <!-- ========================= -->
  <!-- Background -->
  <!-- ========================= -->
  <rect x="0" y="0" width="900" height="440" rx="18" fill="var(--bg)" />

  <!-- ========================= -->
  <!-- Header -->
  <!-- ========================= -->
  <g class="fade-in" style="animation-delay: 0ms">
    <text x="32" y="48" class="title">
      {{.Stats.Username}} · Open Source Stats
    </text>
    <text x="32" y="72" class="subtitle">
      Generated {{.Stats.GeneratedAt}}
    </text>
  </g>

  <!-- ========================= -->
  <!-- Metrics Row -->
  <!-- ========================= -->

  <!-- Projects -->
  <g class="fade-in" style="animation-delay: 100ms">
    <rect x="32" y="96" width="260" height="96" class="card"/>
    <rect x="32" y="96" width="260" height="96" class="card-glass"/>
    <rect x="44" y="96" width="236" height="3" rx="2" class="accent-strip"/>
    <text x="48" y="132" class="metric-label">Projects</text>
    <text x="48" y="166" class="metric-value">{{.TotalProjects}}</text>
  </g>

  <!-- PRs -->
  <g class="fade-in" style="animation-delay: 150ms">
    <rect x="320" y="96" width="260" height="96" class="card"/>
    <rect x="320" y="96" width="260" height="96" class="card-glass"/>
    <rect x="332" y="96" width="236" height="3" rx="2" class="accent-strip"/>
    <text x="336" y="132" class="metric-label">PRs Merged</text>
    <text x="336" y="166" class="metric-value">{{.TotalPRs}}</text>
  </g>

  <!-- Lines Changed -->
  <g class="fade-in" style="animation-delay: 200ms">
    <rect x="608" y="96" width="260" height="96" class="card"/>
    <rect x="608" y="96" width="260" height="96" class="card-glass"/>
    <rect x="620" y="96" width="236" height="3" rx="2" class="accent-strip"/>
    <text x="624" y="132" class="metric-label">Lines Changed</text>
    <text x="624" y="166" class="metric-value">{{.TotalLines}}</text>
  </g>

  <!-- ========================= -->
  <!-- Top Contributions -->
  <!-- ========================= -->
  <g class="fade-in" style="animation-delay: 260ms">
    <text x="32" y="224" class="subtitle">
      Top Contributions
    </text>
  </g>

  {{range $i, $r := .TopContributions}}
  <g class="fade-in" style="animation-delay: {{add 300 (mul $i 100)}}ms">
    <rect x="{{add 32 (mul $i 288)}}" y="240" width="260" height="108" class="card"/>
    <rect x="{{add 32 (mul $i 288)}}" y="240" width="260" height="108" class="card-glass"/>
    <rect x="{{add 44 (mul $i 288)}}" y="240" width="236" height="3" rx="2" class="accent-strip"/>

    <text x="{{add 48 (mul $i 288)}}" y="276" class="repo-name">
      {{$r.RepoName}}
    </text>
    <text x="{{add 48 (mul $i 288)}}" y="302" class="repo-meta">
      ⭐ {{$r.Stars}} · {{$r.PRs}} PRs
    </text>
  </g>
  {{end}}

</svg>
`
