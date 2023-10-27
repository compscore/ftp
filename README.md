# FTP

## Check Format

```yaml
- name:
  release:
    org: compscore
    repo: ftp
    tag: latest
  credentials:
    username:
    password:
  target:
  command:
  expectedOutput:
  weight:
  options:
    exists:
    match:
    substring_match:
    regex_match:
    sha256:
    md5:
    sha1:
```

## Parameters

|     parameter     |            path            |   type   |   default   | required | description                                           |
| :---------------: | :------------------------: | :------: | :---------: | :------: | :---------------------------------------------------- |
|      `name`       |          `.name`           | `string` |    `""`     |  `true`  | `name of check (must be unique)`                      |
|       `org`       |       `.release.org`       | `string` |    `""`     |  `true`  | `organization that check repository belongs to`       |
|      `repo`       |      `.release.repo`       | `string` |    `""`     |  `true`  | `repository of the check`                             |
|       `tag`       |       `.release.tag`       | `string` |  `latest`   | `false`  | `tagged version of check`                             |
|    `username`     |  `.credentials.username`   | `string` | `anonymous` | `false`  | `username of ftp user`                                |
|    `password`     |  `.credentials.password`   | `string` | `anonymous` | `false`  | `default password of ftp user`                        |
|     `target`      |         `.target`          | `string` |    `""`     |  `true`  | `ftp server network location`                         |
|     `command`     |         `.command`         | `string` |    `""`     | `false`  | `file to check against expectedOutput`                |
| `expectedOutput`  |     `.expectedOutput`      | `string` |    `""`     | `false`  | `expected output of file based on options`            |
|     `weight`      |         `.weight`          |  `int`   |     `0`     |  `true`  | `amount of points a successful check is worth`        |
|     `exists`      |     `.options.exists`      |  `bool`  |  `false `   | `false`  | `check targeted file exists and can be accessed`      |
|      `match`      |      `.options.match`      |  `bool`  |   `false`   | `false`  | `check contents of targeted file are exact match`     |
| `substring_match` | `.options.substring_match` |  `bool`  |   `false`   | `false`  | `check contents of targeted file are substring match` |
|   `regex_match`   |   `.options.regex_match`   |  `bool`  |   `false`   | `false`  | `check contents of targeted file are regex match`     |
|     `sha256`      |     `.options.sha256`      |  `bool`  |  `false `   | `false`  | `check sha256 hash of targeted file matches hash`     |
|       `md5`       |       `.options.md5`       |  `bool`  |   `false`   | `false`  | `check md5 hash of targeted file matches hash`        |
|      `sha1`       |      `.options.sha1`       |  `bool`  |   `bool`    | `false`  | `check sha1 hash of targeted file matches hash`       |

## Examples

```yaml
- name: test.rebex.net-ftp
  release:
    org: compscore
    repo: ftp
    tag: latest
  credentials:
    username: demo
    password: password
  target: test.rebex.net
  expectedOutput: b004de45d8a133e9713a369f9c912237e8ad35dd9140c0279d27bada067797f4
  weight: 1
  command: readme.txt
  options:
    exists:
    sha256:
```

```yaml
- name: host_a-ftp
  release:
    org: compscore
    repo: ftp
    tag: latest
  credentials:
    username: john_doe
    password: changeme123!
  target: 10.{ .Team }.1.1
  expectedOutput: ^According to all known laws of aviation
  weight: 1
  command: bee_movie_script.txt
  options:
    regex_match:
```
