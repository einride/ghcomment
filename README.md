ghcomment
=========

A simple CLI tool for commenting on GitHub PRs.

Reasoning
---------

Most tools out there will just post a new comment to a GitHub PR but in heavily developed, long living PRs or when the comments are large, lots of these comments can become a burden to read the PR and also to find out which comments are still relevant and which are now old and outdated.

The main feature `ghcomment` was created for is that it always updates an existing comment, if one already exists previously instead of creating a new one.

Usage
-----

```
Usage of ghcomment:

  -comment string
        A PR Comment
  -pr int
        A Pull Request number
  -repo string
        A GitHub repository on the format <org>/<repository>, or the local git repo if empty
  -signkey string
        A key used to create the signature
  -token string
        A github token
```

SignKey
-------

By using the `-signkey` flag, ghcomment can keep track of several different comments at the same time. The sign key can be any string which is hashed and posted along with the comment. When `ghcomment` runs with a sign key, it will look for a comment in the PR with a hash of that key and if it exist, it updates it instead of creating a new one
