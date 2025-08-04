# Git workflows, submodule, GitHub Actions
The idea is to make a GitHub Actions repo that house all other source code repos as its submodules, and add a few action scripts, webhook which triggers auto-build upon event received.

```
workflows-repo/
├── .github/workflows/
│   ├── auto-build.yml
│   └── manual-build.yml
├── services/
│   ├── service-1/          # submodule
│   ├── service-2/          # submodule
│   └── ...                 # submodule
├── magefile.go             # for build
├── docker-compose.yml      # for local testing
└── scripts/
    └── update-submodules.sh
```

To create a submodule, stand in the `workflows` repo, then run: 
```
git submodule add <url-to-repo>
``` 
where <url-to-repo> can be HTTPS url, or SSH url.
