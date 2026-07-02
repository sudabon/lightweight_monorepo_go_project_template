#!/usr/bin/env sh
set -eu

rule_files="
architecture.md
backend-conventions.md
frontend-conventions.md
git.md
templates.md
testing.md
"

for file in $rule_files; do
  diff -u ".claude/rules/$file" ".codex/rules/$file"
  diff -u ".claude/rules/$file" ".cursor/rules/$file"
done
