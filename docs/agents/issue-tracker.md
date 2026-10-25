# Issue tracker: GitHub

Issues and specs live in BBBuyan/open-monastery on GitHub.
Use the gh CLI from this repo.

## Conventions

- Create: gh issue create --title "..." --body-file <file>
- Read: gh issue view <number> --json number,title,body,labels,comments
- List: gh issue list --state open --json number,title,body,labels
- Comment: gh issue comment <number> --body-file <file>
- Label: gh issue edit <number> --add-label "<label>"
- Remove label: gh issue edit <number> --remove-label "<label>"
- Close: gh issue close <number> --comment "..."

Use a temporary file for multiline bodies.

When a skill says "publish to the issue tracker", create a GitHub issue.
When it says "fetch the relevant ticket", read the GitHub issue.

## Pull requests as a triage surface

PRs as a request surface: no.

## Wayfinding

Keep the map in an issue labelled wayfinder:map.
Link child tickets using GitHub sub-issues; where unavailable, use a
task list in the map and "Part of #<map>" in each child.

Label children wayfinder:<type>, where type is research, prototype,
grilling, or task. Assign claimed tickets to the driving developer.

Record blockers using GitHub issue dependencies; where unavailable,
put "Blocked by: #<number>" at the top of the child body.
A ticket is unblocked when all blockers are closed.

Choose the first open, unassigned, unblocked child in map order.
Claim it with gh issue edit <number> --add-assignee @me.
Resolve it by commenting with the answer, closing the issue, and
adding a summary and link to the map's Decisions-so-far.
