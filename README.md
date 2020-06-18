# Hegemonie

[![CircleCi](https://circleci.com/gh/jfsmig/hegemonie.svg?style=svg)](https://app.circleci.com/pipelines/github/jfsmig/hegemonie)
[![Codecov](https://codecov.io/gh/jfsmig/hegemonie/branch/master/graph/badge.svg)](https://codecov.io/gh/jfsmig/hegemonie)
[![Codacy](https://app.codacy.com/project/badge/Grade/bf7c2872c60445c99f914d31d7b213ae)](https://www.codacy.com/manual/jfsmig/hegemonie?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=jfsmig/hegemonie&amp;utm_campaign=Badge_Grade)
[![MPL-2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

Hegemonie is an online management, strategy and diplomacy RPG. The current
repository is a reboot of what [Hegemonie](http://www.hegemonie.be) was
between 1999 and 2003. It is under heavily inactive construction.

The game engine is a [4X](https://en.wikipedia.org/wiki/4X) technical basis
for larger RPG games. The description of the game will sound familiar to RPG
players:
 * Each **user** manages characters, each **character** manages cities
 * Each **city** produces **resources**, grows a **knowledge** tree, evolves
   with **building** and trains **troops**.
 * **Armies** can be setup to gather some troops and execute orders across
   the map.

## Getting Started

Simply build & install like this:

```
set -e
set -x
BASE=github.com/jfsmig/hegemonie
go get "${BASE}"
go mod download "${BASE}"
go install "${BASE}"
```

For more information, please refer to the page with the [technical elements](./TECH.md).

## How To Contribute

Contributions are what make the open source community such an amazing place.
Any contributions you make are greatly appreciated.

 1. Fork the Project
 2. Create your Feature Branch (git checkout -b feature/AmazingFeature)
 3. Commit your Changes (git commit -m 'Add some AmazingFeature')
 4. Push to the Branch (git push origin feature/AmazingFeature)
 5. Open a Pull Request

## License

Distributed under the MPLv2 License. See [LICENSE](./LICENSE) for more information.

We strongly believe in Open Source for many reasons:
  * For the purpose of a better user experience, because the value of a game
    instance is in its players and in the description of the world. Therefore
    the game engine should focus on allowing a rich world and letting a game
    master to populate instances with an awesome cet of players.
  * For software quality purposes because a software with open sources is the best
    way to have its bugs identified and fixed as soon as possible.
  * For a greater adoption, we chosed a deliberatly liberal license so that
    there cannot be any legal concern.

## Contact

Follow the development on GitHub with the [jfsmig/hegemonie](https://github.com/jfsmig/hegemonie) project.

Follow the community on our Facebook page [Hegemonie.be](https://www.facebook.com/hegemonie.be).

## Acknowledgements

We welcome any volunteer and we already have a list of [amazing authors of Hegemonie](./AUTHORS.md).

