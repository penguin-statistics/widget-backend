<img src="https://penguin.upyun.galvincdn.com/logos/penguin_stats_logo.png"
     alt="Penguin Statistics - Logo"
     width="96px" />

# Penguin Statistics - Widget `backend`
[![Status](https://img.shields.io/badge/status-staging-orange)](#readme)
[![Language](https://img.shields.io/badge/using-Go-%2300add8?logo=go)](#readme)
[![Go Version](https://img.shields.io/github/go-mod/go-version/penguin-statistics/widget-backend)](https://github.com/penguin-statistics/widget-backend/blob/main/go.mod)
[![GoDoc](https://godoc.org/github.com/penguin-statistics/widget-backend?status.svg)](https://godoc.org/github.com/penguin-statistics/widget-backend)
[![Go Report Card](https://goreportcard.com/badge/github.com/penguin-statistics/widget-backend)](https://goreportcard.com/report/github.com/penguin-statistics/widget-backend)
[![License](https://img.shields.io/github/license/penguin-statistics/widget-backend)](https://github.com/penguin-statistics/widget-backend/blob/main/LICENSE)
[![Last Commit](https://img.shields.io/github/last-commit/penguin-statistics/widget-backend)](https://github.com/penguin-statistics/widget-backend/commits/main)

This is the **backend** project repository for the [Penguin Statistics](https://penguin-stats.io/?utm_source=github) widget.

## Important Concepts Outlined
### Why using a different backend rather to just use the [`backend`](https://github.com/penguin-statistics/backend)?
First of all, Go **might** bring us better performance. I used word *might* because there's no benchmark comparison available to give a solid proof, but given by common knowledge Go tends to have better performance. That being said, performance in this particular service is critical since we'd need to probably calculate every request we've got and that request is probably not cacheable or, is not worth to cache given by the short amount of period that the result matrix may update.

Also, adding such feature to the main `backend` project tends to add too much complexity to the project, which may decrease the performance of the project, and more severely after some discussion with the main contributor of such project we concluded that adding such feature to the `backend` project would have a high possibility to impact the readability of the code heavily. Therefore, we have decided to maintain a dedicated project that decouples the logics within that project, where the `backend` project would just to focus on the actual calculation of matrix and to handle other data-sensitive jobs, and this `widget-backend` project would to just focus on manipulating on the already calculated matrices from the `backend` project.

Finally, and probably the most critical reason, is that the widget shall be deployed on a dedicated subdomain, considering healthy status of caches, maintainability and the permanence nature of widget URL. Given by the need to deploy on a dedicated subdomain, a dedicated project would be the most feasible since conditioning the `Host` header inside the same project doesn't always look neat. :D

### Why using `/matrix/:server` as the prefix of all matrix requests, rather to just use a query string?
First of all, sets of records for different servers are by nature completely different sets of resources, and therefore by the semantics of REST API such resources shall have different `path` in order to represent the difference rather to use a query string which is often used to filter records by condition from a defined set of records. 

Moreover, using a completely different path for different resources help with the HTTP-level cache servers to properly cache the request based on the request path rather than the need of parsing query string. Also, to mention that, from experiences we had with Aliyun CDN which we used for our CN mirror site, their cache-matching mechanism tends to either not match query string at all, or to match every single query strings. This leads to a vulnerability where the one could initiate a cache penetration attack towards our endpoint, which they could just use a random query string that changes everytime in order to invalidate the cache that's already stored on the CDN. Given by such experience we've got, we've decided that on the widget backend we would distinguish the server request using path rather to use a query string. This also improves the clarity and readability of the request URL itself.

And finally, given by the possibility (and a very high one) that any time in the future Penguin Statistics project may lose funds which by then, server cost would be infeasible which lead to the possibility of having a proper archive of the project. By using the `path` way to represent all differences across resources, we could ultimately utilize the GitHub Pages where it allows us to render all matrices at once and store them as html pages as format `/matrix/:server/:type/:id.html`. Thereafter, the users of our widget wouldn't need to change anything, including the URL, and would still be able to see the archived widgets.

Given by such reasons we've concluded to use such format by a tradeoff to lose consistencies from our [Public API](https://developer.penguin-stats.io/docs/).

## Maintainers
This project has mainly being maintained by the following contributors (in alphabetical order):
- [GalvinGao](https://github.com/GalvinGao)

> The full list of active contributors of the *Penguin Statistics* project can be found at the [Team Members page](https://penguin-stats.io/about/members) of the website.

## How to contribute?
Our contribute guideline can be found at [Penguin Developers](https://developer.penguin-stats.io). PRs are always more than welcome!