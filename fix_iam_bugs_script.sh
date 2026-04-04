#!/bin/bash
sed -i '' -e 's/- \[ \] \*\*Wrong error code when user is not active\*\*/- \[x\] \*\*Wrong error code when user is not active\*\*/g' /Users/11176728/gemini-cli/dev-turms-im-refactor/turms-go/docs/pending_bugs.md
sed -i '' -e 's/- \[ \] \*\*Missing "deleted" check\*\*/- \[x\] \*\*Missing "deleted" check\*\*/g' /Users/11176728/gemini-cli/dev-turms-im-refactor/turms-go/docs/pending_bugs.md
sed -i '' -e 's/- \[ \] \*\*Password comparison is plain string equality instead of using PasswordManager\*\*/- \[x\] \*\*Password comparison is plain string equality instead of using PasswordManager\*\*/g' /Users/11176728/gemini-cli/dev-turms-im-refactor/turms-go/docs/pending_bugs.md
sed -i '' -e 's/- \[ \] \*\*Granted response returns nil permissions instead of all permissions\*\*/- \[x\] \*\*Granted response returns nil permissions instead of all permissions\*\*/g' /Users/11176728/gemini-cli/dev-turms-im-refactor/turms-go/docs/pending_bugs.md
sed -i '' -e 's/- \[ \] \*\*Finds full user record instead of separate targeted queries\*\*/- \[x\] \*\*Finds full user record instead of separate targeted queries\*\*/g' /Users/11176728/gemini-cli/dev-turms-im-refactor/turms-go/docs/pending_bugs.md

