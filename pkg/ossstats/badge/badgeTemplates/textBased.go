package badgetemplates

const TextBasedCompact = `<svg
  width="280"
  height="32"
  viewBox="0 0 280 32"
  xmlns="http://www.w3.org/2000/svg"
  role="img"
  aria-label="OSS Contributions">
  <defs>
    <style>
      .bg { fill: {{.Colors.Background}}; }
      .content {
        font-family: system-ui, -apple-system, BlinkMacSystemFont,
                     "Segoe UI", Helvetica, Arial, sans-serif;
        font-size: 11px;
        fill: {{.Colors.Text}};
      }
      .bold { font-weight: 700; }
      .sep { fill: {{.Colors.TextSecondary}}; }
    </style>
  </defs>
  <rect class="bg" width="280" height="32" rx="16"/>
  <text class="content" x="140" y="19.5" text-anchor="middle">
    <tspan class="bold">OSS</tspan>
    <tspan class="sep"> · </tspan>
    <tspan class="bold">{{.TotalProjects}}</tspan>
    <tspan> projects</tspan>
    <tspan class="sep"> · </tspan>
    <tspan class="bold">{{.TotalPRs}}</tspan>
    <tspan> PRs Merged</tspan>
  </text>
</svg>`

// TextBasedSummary is the SVG template for the Summary badge style (400x200)
const TextBasedSummary = `<svg
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
      .username {
        font-size: 18px;
        font-weight: 700;
        fill: {{.Colors.Text}};
      }
      .subtitle {
        font-size: 11px;
        fill: {{.Colors.TextSecondary}};
      }
      .stat-value {
        font-size: 28px;
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
  <text class="username" x="28" y="45">@{{.Stats.Username}}</text>
  <text class="subtitle" x="28" y="62">Open Source Contributions</text>
  <!-- Stats -->
  <text class="stat-value" x="60" y="127" text-anchor="middle">{{.TotalProjects}}</text>
  <text class="stat-label" x="60" y="141" text-anchor="middle">Projects</text>
  <text class="stat-value" x="180" y="127" text-anchor="middle">{{.TotalPRs}}</text>
  <text class="stat-label" x="180" y="141" text-anchor="middle">PRs Merged</text>
  <text class="stat-value" x="300" y="127" text-anchor="middle">{{.TotalLines}}</text>
  <text class="stat-label" x="300" y="141" text-anchor="middle">Lines Changed</text>
</svg>`

// textBasedDetailedTemplate is the SVG template for the Detailed badge style (400x320)
const TextBasedDetailed = `
	{{ $SVGHeight := add 300 (mul 56 (len .TopContributions)) }}
	<svg
  width="720"
  height="{{$SVGHeight}}"
  viewBox="0 0 720 {{$SVGHeight}}"
  xmlns="http://www.w3.org/2000/svg"
  role="img"
  aria-label="Open Source Contributions">

  <defs>
    <style>
      text {
        font-family: system-ui, -apple-system, BlinkMacSystemFont,
                     "Segoe UI", Helvetica, Arial, sans-serif;
      }

      .bg {
        fill: {{.Colors.Background}};
      }

      .username {
        font-size: 30px;
        font-weight: 900;
        fill: {{.Colors.Text}};
        letter-spacing: -0.6px;
      }

      .subtitle {
        font-size: 13px;
        fill: {{.Colors.TextSecondary}};
      }

      .stat {
        font-size: 26px;
        font-weight: 800;
        fill: {{.Colors.Text}};
      }

      .stat-label {
        font-size: 12px;
        fill: {{.Colors.TextSecondary}};
      }

      .divider {
        stroke: {{.Colors.Border}};
        stroke-width: 1;
      }

      .section-title {
        font-size: 13px;
        font-weight: 700;
        fill: {{.Colors.TextSecondary}};
        letter-spacing: 0.1em;
      }

      .repo-name {
        font-size: 16px;
        font-weight: 700;
        fill: {{.Colors.Text}};
      }

      .repo-stars {
        font-size: 13px;
        font-weight: 700;
        fill: {{.Colors.Accent}};
      }

      .repo-meta {
        font-size: 12px;
        fill: {{.Colors.TextSecondary}};
      }
    </style>
  </defs>

  <!-- Background -->
  <rect
    class="bg"
    x="0"
    y="0"
    width="720"
    height="{{$SVGHeight}}"
    rx="20"/>

  <!-- Header -->
  <text class="username" x="48" y="68">
    {{.Stats.Username}}
  </text>

  <text class="subtitle" x="48" y="92">
    Open source contributions
  </text>

  <!-- Stats -->
  <g transform="translate(48, 132)">
    <text class="stat">{{.TotalProjects}}</text>
    <text class="stat-label" y="22">Projects</text>
  </g>

  <g transform="translate(220, 132)">
    <text class="stat">{{.TotalPRs}}</text>
    <text class="stat-label" y="22">PRs merged</text>
  </g>

  <g transform="translate(390, 132)">
    <text class="stat">{{.TotalLines}}</text>
    <text class="stat-label" y="22">Lines changed</text>
  </g>

  <!-- Divider -->
  <line
    class="divider"
    x1="48"
    y1="196"
    x2="672"
    y2="196"/>

  <!-- Repo Section -->
  <text class="section-title" x="48" y="228">
    TOP REPOSITORIES
  </text>

  {{range $i, $r := .TopContributions}}
  <g transform="translate(48, {{add 260 (mul $i 56)}})">
    <text class="repo-name">
      {{$r.RepoName}}
      <tspan class="repo-stars"> ★ {{$r.Stars}}</tspan>
    </text>

    <text class="repo-meta" y="22">
      {{$r.PRs}} PRs merged
    </text>
  </g>
  {{end}}

</svg>

`
