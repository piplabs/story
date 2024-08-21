# Contributing to Story
Thank you for considering making contributions to Story and related repositories! Below you can find what's the best way that you can help.

## Coding Guideline

[Go coding guidelines](https://github.com/golang/go/wiki/CodeReviewComments)

[Good commit messages](https://chris.beams.io/posts/git-commit/)

## Bug Reports

* Ensure your issue [has not already been reported][1]. It may already be fixed!
* Include the steps you carried out to produce the problem.
* Include the behavior you observed along with the behavior you expected, and
  why you expected it.
* Include any relevant stack traces or debugging output.

## Feature Requests

We welcome feedback with or without pull requests. If you have an idea for how
to improve the project, great! All we ask is that you take the time to write a
clear and concise explanation of what need you are trying to solve. If you have
thoughts on _how_ it can be solved, include those too!

The best way to see a feature added, however, is to submit a pull request.

## Pull Requests

* Before creating your pull request, it's usually worth asking if the code
  you're planning on writing will actually be considered for merging. You can
  do this by [opening an issue][1] and asking. It may also help give the
  maintainers context for when the time comes to review your code.

* Ensure your [commit messages are well-written][2]. This can double as your
  pull request message, so it pays to take the time to write a clear message.

* Add tests for your feature. You should be able to look at other tests for
  examples. If you're unsure, don't hesitate to [open an issue][1] and ask!

* Submit your pull request!
    - Fork the repository on GitHub.
    - Make your changes on your fork repository.
    - Submit a PR.

* Each PR needs to link to a github issue. Please open a github issue first to describe the code changes you plan to make and the problem it solves. Then link the issue in your PR.

* PR title and body follows [conventional commit][3].
   - Title template: `type(app/pkg): concise description`
   - See [allowed types][4]
   - Description must be concise: lower case, no punctuation, no more than 50 characters.
   - Scope must be concise: only a one or two folders; e.g. 'client/cmd' or 'github' or '*'

## Find something to work on

We are always in need of help, be it fixing documentation, reporting bugs or writing some code.
Look at places where you feel best coding practices aren't followed, code refactoring is needed or tests are missing.

If you have questions about the development process,
feel free to [file an issue](https://github.com/piplabs/story/issues/new).

## Code Review

To make it easier for your PR to receive reviews, consider the reviewers will need you to:

* follow [good coding guidelines](https://github.com/golang/go/wiki/CodeReviewComments).
* write [good commit messages](https://chris.beams.io/posts/git-commit/).
* break large changes into a logical series of smaller patches which individually make easily understandable changes, and in aggregate solve a broader issue.

[1]: https://github.com/piplabs/story/issues
[2]: https://chris.beams.io/posts/git-commit/#seven-rules
[3]: https://www.conventionalcommits.org/en/v1.0.0
[4]: https://github.com/conventional-changelog/commitlint/tree/master/%40commitlint/config-conventional#type-enum
