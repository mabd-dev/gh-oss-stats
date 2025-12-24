package badge

// summaryTemplate is the SVG template for the Summary badge style (400x200)
const summaryTemplate = `<svg width="400" height="200" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <style>
      .bg { fill: {{.Colors.Background}}; }
      .border { stroke: {{.Colors.Border}}; stroke-width: 1; fill: none; }
      .text { fill: {{.Colors.Text}}; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif; }
      .title { font-size: 16px; font-weight: 600; }
      .value { font-size: 24px; font-weight: 700; fill: {{.Colors.Accent}}; }
      .label { font-size: 12px; fill: {{.Colors.TextSecondary}}; }
    </style>
  </defs>
  
  <!-- Background -->
  <rect class="bg border" x="0.5" y="0.5" width="399" height="199" rx="6"/>
  
  <!-- Header -->
  <text class="text title" x="20" y="30">Open Source Contributions</text>
  <text class="label" x="20" y="50">@{{.Stats.Username}}</text>
  
  <!-- Stats Grid -->
  <!-- Row 1: Projects and PRs -->
  <g transform="translate(20, 80)">
    <text class="value" x="0" y="0">{{.TotalProjects}}</text>
    <text class="label" x="0" y="18">Projects</text>
  </g>
  
  <g transform="translate(220, 80)">
    <text class="value" x="0" y="0">{{.TotalPRs}}</text>
    <text class="label" x="0" y="18">PRs Merged</text>
  </g>
  
  <!-- Row 2: Commits and Lines -->
  <g transform="translate(20, 140)">
    <text class="value" x="0" y="0">{{.TotalCommits}}</text>
    <text class="label" x="0" y="18">Commits</text>
  </g>
  
  <g transform="translate(220, 140)">
    <text class="value" x="0" y="0">{{.TotalLines}}</text>
    <text class="label" x="0" y="18">Lines Changed</text>
  </g>
</svg>`

// compactTemplate is the SVG template for the Compact badge style (280x28) - Shields.io style
const compactTemplate = `<svg width="280" height="28" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <style>
      .bg-left { fill: {{.Colors.TextSecondary}}; }
      .bg-right { fill: {{.Colors.Accent}}; }
      .text { fill: #ffffff; font-family: Verdana, Geneva, DejaVu Sans, sans-serif; font-size: 11px; font-weight: 600; }
    </style>
  </defs>
  
  <!-- Left section (label) -->
  <rect class="bg-left" x="0" y="0" width="160" height="28" rx="3"/>
  
  <!-- Right section (value) -->
  <rect class="bg-right" x="160" y="0" width="120" height="28" rx="3"/>
  <rect class="bg-right" x="160" y="0" width="3" height="28"/>
  
  <!-- Text -->
  <text class="text" x="80" y="18" text-anchor="middle">OSS Contributions</text>
  <text class="text" x="220" y="18" text-anchor="middle">{{.CompactText}}</text>
</svg>`

// minimalTemplate is the SVG template for the Minimal badge style (120x28)
const minimalTemplate = `<svg width="120" height="28" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <style>
      .bg { fill: {{.Colors.Accent}}; }
      .text { fill: #ffffff; font-family: Verdana, Geneva, DejaVu Sans, sans-serif; font-size: 11px; font-weight: 600; }
    </style>
  </defs>
  
  <!-- Background -->
  <rect class="bg" x="0" y="0" width="120" height="28" rx="3"/>
  
  <!-- Text -->
  <text class="text" x="60" y="18" text-anchor="middle">{{.MinimalText}}</text>
</svg>`

// detailedTemplate is the SVG template for the Detailed badge style (400x320)
const detailedTemplate = `<svg width="400" height="320" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <style>
      .bg { fill: {{.Colors.Background}}; }
      .border { stroke: {{.Colors.Border}}; stroke-width: 1; fill: none; }
      .text { fill: {{.Colors.Text}}; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif; }
      .title { font-size: 16px; font-weight: 600; }
      .value { font-size: 24px; font-weight: 700; fill: {{.Colors.Accent}}; }
      .label { font-size: 12px; fill: {{.Colors.TextSecondary}}; }
      .repo-name { font-size: 13px; font-weight: 500; }
      .repo-stat { font-size: 11px; fill: {{.Colors.TextSecondary}}; }
      .divider { stroke: {{.Colors.Border}}; stroke-width: 1; }
    </style>
  </defs>
  
  <!-- Background -->
  <rect class="bg border" x="0.5" y="0.5" width="399" height="319" rx="6"/>
  
  <!-- Header -->
  <text class="text title" x="20" y="30">Open Source Contributions</text>
  <text class="label" x="20" y="50">@{{.Stats.Username}}</text>
  
  <!-- Summary Stats -->
  <g transform="translate(20, 75)">
    <text class="value" x="0" y="0">{{.TotalProjects}}</text>
    <text class="label" x="0" y="18">Projects</text>
  </g>
  
  <g transform="translate(120, 75)">
    <text class="value" x="0" y="0">{{.TotalPRs}}</text>
    <text class="label" x="0" y="18">PRs</text>
  </g>
  
  <g transform="translate(220, 75)">
    <text class="value" x="0" y="0">{{.TotalCommits}}</text>
    <text class="label" x="0" y="18">Commits</text>
  </g>
  
  <!-- Divider -->
  <line class="divider" x1="20" y1="115" x2="380" y2="115"/>
  
  <!-- Top Contributions Header -->
  <text class="label" x="20" y="135">Top Contributions</text>
  
  <!-- Top Contributions List -->
  {{range $i, $c := .TopContributions}}
  <g transform="translate(20, {{add 155 (mul $i 30)}})">
    <text class="text repo-name" x="0" y="0">{{$c.RepoName}}</text>
    <text class="repo-stat" x="0" y="15">⭐ {{$c.Stars}}  •  {{$c.PRs}} PRs</text>
  </g>
  {{end}}
</svg>`
