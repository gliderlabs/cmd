// magical github automation
//
// add hooks similar to below based on:
// https://developer.github.com/v3/activity/events/types/
//
// use github object according to docs:
// http://mikedeboer.github.io/node-github/

const Promise = require("bluebird");

const pleaseRebase = "Please rebase your PR, thanks!";

function onPush(github, event, cb) {
  console.log("onPush");
  if (event.ref != "refs/heads/master") {
    console.log("Not master");
    cb();
  }
  Promise.promisifyAll(github);
  github.pullRequests.getAll({
    user: event.repository.owner.name,
    repo: event.repository.name,
    base: "master"

  }).then(function(prs) {
    if (prs.length == 0) {
      console.log("No PRs");
      return cb();
    }
    Promise.all(
      prs.map(function(pr) {
        return github.pullRequests.get({
          user: pr.base.repo.owner.login,
          repo: pr.base.repo.name,
          number: pr.number
        })
      })
    ).then(function(prs) {
      Promise.all(
        prs.map(function(pr) {
          if (pr.mergeable == false) {
            return github.issues.getComments({
              user: pr.base.repo.owner.login,
              repo: pr.base.repo.name,
              number: pr.number,
              per_page: 100
            }).then(function(comments) {
              if (comments.length > 0 && comments[comments.length-1].body == pleaseRebase) {
                return Promise.resolve();
              }
              return github.issues.createComment({
                user: pr.base.repo.owner.login,
                repo: pr.base.repo.name,
                number: pr.number,
                body: pleaseRebase
              });
            });
          } else {
            return Promise.resolve();
          }
        })
      ).then(function() {
        console.log("done!");
        cb();
      });
    });
  }).catch(function(err) {
    if (err) {
      console.log("error:", err.message);
      return cb();
    }
  });
}

function onIssueComment(github, event, cb) {
    const {number, state} = event.issue;
    const owner = event.repository.owner.login;
    const repo = event.repository.name;
    const author = event.comment.user.login;
    const msg = event.comment.body;

    if (msg.toLowerCase() !== 'lgtm') {
        console.log("message isn't lgtm");
        return cb();
    }
    if (state !== 'open') {
        console.log("issue isn't open");
        return cb();
    }

    github.repos.checkCollaborator({
        user: owner,
        repo: repo,
        collabuser: author
    }).then(
        res => {
            // Atempt to merge PR.
            return github.pullRequests.merge({
                user: owner,
                repo: repo,
                number: number,
                squash: true
            });
        },
        err => {
            console.log(JSON.stringify({
                "error": err
            }));
            throw new Error("author is not collaborator");
        }
    ).then(
        res => {
            console.log(JSON.stringify(res))
        },
        err => {
            console.log(JSON.stringify({
                "error": err
            }));
            if (err.code == 405) {
                throw new Error("Pull Request is not mergeable");
            }
            if (err.code == 404) {
                throw new Error("PR for issue number not found");
            }
        }
    ).then(cb).catch(e => {
        console.log(e);
        cb();
    });
}

function onPing(github, event, cb) {
  console.log("ping");
}
