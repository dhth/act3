import subprocess

repos = [
    "dhth/commits",
    "dhth/cueitup",
    "dhth/dstll",
    "dhth/ecsv",
    "dhth/hours",
    "dhth/kplay",
    "dhth/mult",
    "dhth/outtasync",
    "dhth/prs",
    "dhth/punchout",
    "dhth/schemas",
]

actions = [
    "build",
    "vulncheck",
    "release",
]

for repo in repos:
    for action in actions:
        command = [
            "gh",
            "api",
            f"repos/{repo}/actions/workflows/{action}.yml",
            "--jq",
            ".node_id",
        ]
        result = subprocess.run(command, capture_output=True, text=True)
        if result.returncode == 0:
            print(
                f"""- id: {result.stdout.strip()}
  repo: {repo}
  name: {action}
"""
            )
