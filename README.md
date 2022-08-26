# Smers
### Application for smer notes. Technologies Used
- React
- Typescript
- Golang
- PostgreSQL

### Start client dev (web folder)
- `yarn`
- `yarn dev`

### Start backend dev (backend folder)
* Installing packages (folder backend/app)
- `go install backend/internal/app backend/internal/config backend/pkg/logging`

* Running backend
- `go run ./app/cmd/server/main.go` — run
- `goose create add_some_column sql` — create migration
- `goose up` — apply migrations

#### Examples:
`goose postgres "user=postgres dbname=stack sslmode=disable" up`

### Flow
- New branch (#issue)</li>
- `git rebase -i @~n` - сливаем в один комит (1 коммит на фичу)
- PR в develop
- Code review
- Merge

### Semantic commit messages
- _**feat**_: (new feature for the user, not a new feature for build script)
- _**fix**_: (bug fix for the user, not a fix to a build script)
- _**chore**_: (updating grunt tasks etc; no production code change)
- _**build:**_: (update for build script)
- _**refactor**_: (refactoring production code, eg. renaming a variable)
- _**redesign**_: (changes related to redesign of interface)
- _**test**_: (adding missing tests, refactoring tests; no production code change)
- _**perf**_: (changes related to performance)
- _**docs**_: (changes to the documentation)
- _**style**_: (formatting, missing semi colons, etc; no production code change)

#### Examples:
`feat: ID column for smers table`
`refactor: smers dialog`
`style: formatted by eslint`

### Project Status
In Progress

### Contacts

<a href="https://github.com/Videot4pe"><img target="_blank" src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/github/github-original.svg" style="width: 10%;"></a>
<a href="https://github.com/stroooooom"><img target="_blank" src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/github/github-original.svg" style="width: 10%;"></a>
