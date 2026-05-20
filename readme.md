# Bookings and Reservations

The repository for [Building Modern Web Applications with Go](https://www.udemy.com/course/building-modern-web-applications-with-go/?referralCode=0415FB906223F10C6800).

## GitKraken / file watching

This folder is the **git repository root**. If you use GitKraken, open this directory (not the parent `goProjects` folder).

- Open repo: `./scripts/open-git-repo.sh`
- Fix Linux watcher limits: `./scripts/setup-inotify.sh` then restart GitKraken



- Built in Go version 1.15
- Uses the [chi router](github.com/go-chi/chi)
- Uses [alex edwards scs session management](github.com/alexedwards/scs)
- Uses [nosurf](github.com/justinas/nosurf)