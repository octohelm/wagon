# Wagon

## Install

```shell
curl -sSLf https://raw.githubusercontent.com/octohelm/wagon/main/install.sh | sudo sh
```

## Difference from `dagger-cue`

* No `core.#Export`
    * for each action, we could just use `do --output` to export to somewhere:
        * when action contains `output: core.#FS`, will export all fs contents
        * when action contains `output: core.#Image`, will export the oci image
* Not required `dagger.#Plan`
    * actions path just under `actions` will be the task list for `wagon do <action_path>`
