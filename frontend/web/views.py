from pathlib import Path
import os
import shlex
import shutil
import subprocess

from django.shortcuts import render


def _build_cli_command() -> tuple[list[str], Path]:
    repo_root = Path(__file__).resolve().parents[2]
    binary_name = "yt-grab.exe" if os.name == "nt" else "yt-grab"
    local_binary = repo_root / binary_name

    if shutil.which("yt-grab"):
        return ["yt-grab"], repo_root
    if local_binary.exists():
        return [str(local_binary)], repo_root
    return ["go", "run", "./cmd/yt-grab"], repo_root


def _parse_download_input(raw: str) -> tuple[str, bool]:
    tokens = shlex.split(raw.strip())
    if not tokens:
        raise ValueError("Paste a YouTube link.")

    url = tokens[0]
    if not (url.startswith("http://") or url.startswith("https://")):
        raise ValueError("Link must start with http:// or https://")

    audio = False
    for token in tokens[1:]:
        if token == "--audio":
            audio = True
        else:
            raise ValueError("Only --audio is allowed after the link.")

    return url, audio


def home(request):
    context = {"input_value": "", "command": "", "output": "", "error": "", "success": False}

    if request.method == "POST":
        raw = request.POST.get("download_input", "").strip()
        context["input_value"] = raw
        try:
            url, audio = _parse_download_input(raw)
            base_cmd, cwd = _build_cli_command()
            cmd = [*base_cmd, url]
            if audio:
                cmd.append("--audio")

            context["command"] = " ".join(shlex.quote(part) for part in cmd)
            result = subprocess.run(
                cmd,
                cwd=cwd,
                capture_output=True,
                text=True,
                timeout=600,
            )
            combined_output = (result.stdout or "") + (result.stderr or "")
            context["output"] = combined_output.strip() or "No output"
            context["success"] = result.returncode == 0
            if not context["success"] and "yt-dlp not found" in combined_output:
                context["error"] = (
                    "yt-dlp is not installed. On Windows run: winget install yt-dlp.yt-dlp"
                )
            elif not context["success"]:
                context["error"] = "Download failed. See output below."
        except subprocess.TimeoutExpired:
            context["error"] = "Download timed out after 10 minutes."
        except Exception as exc:
            context["error"] = str(exc)

    return render(request, "web/home.html", context)
