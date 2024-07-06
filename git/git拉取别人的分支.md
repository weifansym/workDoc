在Git中拉取别人的分支可以通过以下步骤实现:

1. 首先，使用 `git clone` 命令将远程仓库克隆到本地。例如，如果要克隆名为 `example_repo` 的仓库，可以运行以下命令：
```
git clone https://github.com/username/example_repo.git
```
其中 `username` 是远程仓库所有者的用户名。

2. 接下来，使用 `git branch -r` 命令查看远程仓库中的分支列表。例如，如果要查看所有的远程分支，可以运行以下命令：
```
git branch -r
```
这将显示远程仓库的分支列表。

3. 选择要拉取的分支。例如，如果要拉取名为 `other_branch` 的分支，可以运行以下命令：
```
git checkout -b other_branch origin/other_branch
```
这将在本地创建一个名为 `other_branch` 的分支，并将其与远程分支 `origin/other_branch` 关联。

4. 最后，可以使用 `git pull` 命令将远程分支的更改拉取到本地。例如，如果要拉取 `other_branch` 的更改，可以运行以下命令：
```
git pull origin other_branch
```
这将拉取远程分支 `origin/other_branch` 的更改，并将其合并到当前所在的分支。

这样，你就成功拉取了别人的分支，并将其同步到本地仓库中。记得经常使用 `git pull` 命令来更新分支，以保持与远程仓库的同步。


