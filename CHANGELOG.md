# Change Log

Notable changes up to the last release.

## [unreleased]

### ⛰️  Features

- Added mutex to the `LevelDB` store - ([70a3722](https://github.com/Friends-Of-Noso/NosoGo/commit/70a372286b9851a6fccf1b6f8cebac02637dbde7))
- Using `GetMyIP()` - ([9a68755](https://github.com/Friends-Of-Noso/NosoGo/commit/9a68755b65f99b5f5d54a885bfd8a290775b580d))
- Adding `getmyip` - ([24ddea4](https://github.com/Friends-Of-Noso/NosoGo/commit/24ddea4f7cb7afd4f823398321837520e4ea3894))
- Startup checks are done - ([cd695b8](https://github.com/Friends-Of-Noso/NosoGo/commit/cd695b84290e7cc34ec06bf90a6766c96fc9c609))

### 🐛 Bug Fixes

- *(cliff.toml)* Removed problematic `typo` util - ([14006cc](https://github.com/Friends-Of-Noso/NosoGo/commit/14006ccc26008069ad42d31825ab4cf4ea0a72a3))
- Proper use of `quit` channel - ([db490c9](https://github.com/Friends-Of-Noso/NosoGo/commit/db490c94cbbe3064f3049842d6d66f6097b8ba1b))
- Finally cracked flag<->config sync - ([76545a8](https://github.com/Friends-Of-Noso/NosoGo/commit/76545a8af53090896e706c23eab3c3848d4d71d9))
- Exit strategies: Signal and `quit` channel - ([5f42ff8](https://github.com/Friends-Of-Noso/NosoGo/commit/5f42ff84d35c939fddeea24ded12ac150a791227))

### 🚜 Refactor

- Now using `PeerList` on `/v1/dns` - ([43aedfd](https://github.com/Friends-Of-Noso/NosoGo/commit/43aedfdff1d8485dcd72aba47ed17bf32dbbeb4d))
- Moved `DNS` resolve helpers to `utils` - ([31f8dcb](https://github.com/Friends-Of-Noso/NosoGo/commit/31f8dcb32a4178b7805e0249e78c08c839f30fe5))
- Network message names - ([4f49c67](https://github.com/Friends-Of-Noso/NosoGo/commit/4f49c6743ed2bc387e90e6d1488e3826b0c60560))

### 📚 Documentation

- *(CHANGELOG)* Removing failed release - ([302118c](https://github.com/Friends-Of-Noso/NosoGo/commit/302118ce222bb4cc8d24fb8a2e2ee961194babee))

### 🎨 Styling

- Better log generation - ([545407e](https://github.com/Friends-Of-Noso/NosoGo/commit/545407e991ea8c048107fedee4336df897ca2241))
- Using `map.Clone()` - ([36f8a90](https://github.com/Friends-Of-Noso/NosoGo/commit/36f8a906485768b185e809c4dc239db8c09b9b7e))
- Better formatting strings - ([21d7ab6](https://github.com/Friends-Of-Noso/NosoGo/commit/21d7ab601870d647041a383f42bbe9b8a48a05ec))
- Using `any` instead of `interface{}` - ([b73008b](https://github.com/Friends-Of-Noso/NosoGo/commit/b73008b49d642a43be7e251d11aa7838ebe4255e))
- Removing commented code - ([2174d5c](https://github.com/Friends-Of-Noso/NosoGo/commit/2174d5c42280e87b93359b0ec7569bd8250069a4))

### 🧪 Testing

- Adding tests for `PeerList` - ([f4db85a](https://github.com/Friends-Of-Noso/NosoGo/commit/f4db85aaabdb0fa568cb9b3eb2c10250e32e7110))
- Using `assert` - ([7c7d501](https://github.com/Friends-Of-Noso/NosoGo/commit/7c7d501647ec7f59fe1efc7bf92ca2a8da366371))
- Skipping `TestBlocksStorageListIncorrect` - ([de54cd2](https://github.com/Friends-Of-Noso/NosoGo/commit/de54cd2868df85980ab064fde64e9edb0d51b9bd))
- Deleting via `LevelDB`, not the `OS` - ([9a7a9ab](https://github.com/Friends-Of-Noso/NosoGo/commit/9a7a9ab6a8f48f8d647a44f3fef380053c9b9376))
- WHat the error on `GitHub`? - ([50e5288](https://github.com/Friends-Of-Noso/NosoGo/commit/50e528856ecc340f317bcaa63fce02c7f2c12065))
- Fix not deleting the data folder - ([818640c](https://github.com/Friends-Of-Noso/NosoGo/commit/818640c16f372550c156642e5035597f4722cdbd))

### ⚙️ Miscellaneous Tasks

- Adding `AUTHORS` file - ([36d2f2f](https://github.com/Friends-Of-Noso/NosoGo/commit/36d2f2fcf93a8cb49a8e1e8bf516ed594fb9db3f))
- Fixed the quotes - ([600e7b2](https://github.com/Friends-Of-Noso/NosoGo/commit/600e7b23a933741cb722924b17ef48948b49eacb))
- Using `cut` to get the correct string - ([455b4a3](https://github.com/Friends-Of-Noso/NosoGo/commit/455b4a3ad4581ef70771d8c554bc481690f4e3fd))
- Using proper Go version instad of `stable` - ([95c1311](https://github.com/Friends-Of-Noso/NosoGo/commit/95c131191deebd75a13516b4aad8cbbe47087ed9))
- Only keeping `go-stable` verions - ([02b206f](https://github.com/Friends-Of-Noso/NosoGo/commit/02b206ff8b3c9db78313b8dbcb2b31518f39e543))
- Bumping version after last release - ([98d964d](https://github.com/Friends-Of-Noso/NosoGo/commit/98d964d768f999b7a731cc625d8a8b84cd18031e))

### ◀️ Revert

- Not deleting anything! Windows :facepalm: - ([acd4f65](https://github.com/Friends-Of-Noso/NosoGo/commit/acd4f65ddc0c0b02f8d54aa3801a7615d693e4cd))

## Contributors

* [@gcarreno](https://github.com/gcarreno)

## [0.0.4](https://github.com/Friends-Of-Noso/NosoGo/compare/v0.0.3..v0.0.4) - 2025-07-07

### ⛰️  Features

- NewBlock, NewTransactions dumped into db - ([9d77245](https://github.com/Friends-Of-Noso/NosoGo/commit/9d77245593ab0c2d87f36c17ab158c31ae5393ee))
- Add generics version to store on `LevelDB` - ([786f209](https://github.com/Friends-Of-Noso/NosoGo/commit/786f209e6b58f20e9d4fff2be7c88b582f74ac6e))

### 🚜 Refactor

- Improved `LevelDB` generics store+tests - ([fde453a](https://github.com/Friends-Of-Noso/NosoGo/commit/fde453a583098328538c780e682be976177885eb))

### 🎨 Styling

- Moving dev only code to it's own file - ([09bb4ba](https://github.com/Friends-Of-Noso/NosoGo/commit/09bb4ba88002bb382f9c96e85b96edeab2f8cdef))

### ⚙️ Miscellaneous Tasks

- *(cliff.toml)* New version of the file - ([6137d33](https://github.com/Friends-Of-Noso/NosoGo/commit/6137d33517e9364b4a61a8501195bfedfbbc6518))
- *(main.yaml)* Fixed repeated params - ([72cdaa9](https://github.com/Friends-Of-Noso/NosoGo/commit/72cdaa9617202d80a1ccd29bdc2fad3b656555b5))
- *(main.yaml)* Using cliff action @ main - ([c54f1a1](https://github.com/Friends-Of-Noso/NosoGo/commit/c54f1a1b32eb061b599910e2375f08d4af7bce63))

## Contributors

* [@gcarreno](https://github.com/gcarreno)

## [0.0.3](https://github.com/Friends-Of-Noso/NosoGo/compare/v0.0.2..v0.0.3) - 2025-07-06

### ⛰️  Features

- Node in mode `dns` starts web server - ([3d78ee0](https://github.com/Friends-Of-Noso/NosoGo/commit/3d78ee08e1c94c8cd16e38ea80917efee9ab6df4))
- Legacy address generation and validation - ([29c1f26](https://github.com/Friends-Of-Noso/NosoGo/commit/29c1f263c859256b24de35431179d248d690e277))

### 🐛 Bug Fixes

- Initial messages should be sent to `STDERR` - ([b2a6d56](https://github.com/Friends-Of-Noso/NosoGo/commit/b2a6d56ac88592bd4b6f8ee2d3019f032a991f80))

### 📚 Documentation

- *(README)* Formatting of the note - ([66671a9](https://github.com/Friends-Of-Noso/NosoGo/commit/66671a9a51366e894d0e35d06fcfc79847001ee1))
- *(README)* Note about the status of the project - ([df3f385](https://github.com/Friends-Of-Noso/NosoGo/commit/df3f38587fb98dc561ec6882d3528de1aeca64f8))

### ⚙️ Miscellaneous Tasks

- *(go.mod)* Upgraded `mapstructure` to v2.3.0 - ([6f6355c](https://github.com/Friends-Of-Noso/NosoGo/commit/6f6355c63e9ee664e829829e0d0f1fc99bc97f12))
- Using Go v1.23 for release - ([eeddc4e](https://github.com/Friends-Of-Noso/NosoGo/commit/eeddc4e305a6c2bf3ed36198f8753e58318c6df6))

### Security

- Upgrades mods to latest version - ([b3eb0c0](https://github.com/Friends-Of-Noso/NosoGo/commit/b3eb0c0c2b8dd95a65e0edd867ebaf12472c6616))
- Mapstructure version update - ([1f86570](https://github.com/Friends-Of-Noso/NosoGo/commit/1f86570a767df6370bb82e8be526add3db2d43f6))
- Updating package versions - ([90e28c0](https://github.com/Friends-Of-Noso/NosoGo/commit/90e28c04e9312db3fa2728dc85ec9761b8615c6d))

## Contributors

* [@gcarreno](https://github.com/gcarreno)

## [0.0.2](https://github.com/Friends-Of-Noso/NosoGo/compare/v0.0.1..v0.0.2) - 2024-11-16

### ⚙️ Miscellaneous Tasks

- Bumping version for test of `CD` - ([ab76a08](https://github.com/Friends-Of-Noso/NosoGo/commit/ab76a080822114a8065a2d2e50daea7f6fd752d3))
- Fix missing version on the released files - ([28d8734](https://github.com/Friends-Of-Noso/NosoGo/commit/28d8734068acd9521e970d039eba5cb9c28a986f))
- Bumping version on release action - ([1736d21](https://github.com/Friends-Of-Noso/NosoGo/commit/1736d218915dfe1ac367a842438c89cdf163f3f9))

## Contributors

* [@gcarreno](https://github.com/gcarreno)

## [0.0.1] - 2024-11-16

### ⛰️  Features

- More comms messages - ([736bd80](https://github.com/Friends-Of-Noso/NosoGo/commit/736bd801b45e2f0a88631f5eb3136511edd87d66))
- `API` related constants and `structs` - ([2c116ff](https://github.com/Friends-Of-Noso/NosoGo/commit/2c116ff1ee0e3f87e90ce76bbf4f3e8e0d5f72e1))
- Scaffolding for Blocks and Blocks Status on client - ([3cbde44](https://github.com/Friends-Of-Noso/NosoGo/commit/3cbde44fdbddd6898795be787c631a89a6cc94ea))
- `CI`/CD` and `API` - ([8b0fdca](https://github.com/Friends-Of-Noso/NosoGo/commit/8b0fdcaccc7a204c49b3e49864bf7f3a9c96090b))
- Scaffolding for `nosogocli` - ([88be101](https://github.com/Friends-Of-Noso/NosoGo/commit/88be101c9dbf9e8ebe97615cb54e73fff70e3ab4))
- Crude block propagation - ([6d667e3](https://github.com/Friends-Of-Noso/NosoGo/commit/6d667e317b8570ced8347718a4e7cbc3c390a8d0))
- Adding address and port from config - ([4e07c66](https://github.com/Friends-Of-Noso/NosoGo/commit/4e07c66b1cc8e27e4fc7b1a98d976030523c0c3c))
- Adding `protobuf` network messages - ([f723139](https://github.com/Friends-Of-Noso/NosoGo/commit/f723139563ba6de25ab428298de3d984826d5ffa))
- Node done with blocking and signals - ([1f4f993](https://github.com/Friends-Of-Noso/NosoGo/commit/1f4f993ebcd1aa1c007e854741c278d8b2526068))
- Adding `init` command's code - ([a0843d6](https://github.com/Friends-Of-Noso/NosoGo/commit/a0843d61a8194e3d5ad04e83a8d3739e9645e68d))
- Scaffolding for Network Status - ([ccaab39](https://github.com/Friends-Of-Noso/NosoGo/commit/ccaab392f406753a2266d69907936dbcb25ad737))
- Scaffolding for network messages - ([bd604b0](https://github.com/Friends-Of-Noso/NosoGo/commit/bd604b0b68f7c2029abf156c3c1aace98dc5dfd2))
- Scaffolding for the blocks container - ([7809185](https://github.com/Friends-Of-Noso/NosoGo/commit/7809185bf76772cb93e80fb1d13b0bb03812721d))
- Capacity capped at 255 - ([3a099ea](https://github.com/Friends-Of-Noso/NosoGo/commit/3a099ea6807a4e064978660b9c0177bca30a5331))
- Adding Pascal Short String struct - ([dd80360](https://github.com/Friends-Of-Noso/NosoGo/commit/dd803602b0c44242a44d4035f78402b53aa1b277))
- Adding node command - ([46519de](https://github.com/Friends-Of-Noso/NosoGo/commit/46519de17429ebb139dc123da4f03a8799f39705))
- Adding some default commands - ([37e7183](https://github.com/Friends-Of-Noso/NosoGo/commit/37e7183daf2ab04eede3cb3daf55545a295405e2))

### 🐛 Bug Fixes

- Length and Capacity - ([023d56b](https://github.com/Friends-Of-Noso/NosoGo/commit/023d56b2d73c019afb7f2c683d486b83313a1583))

### 🚜 Refactor

- Using one of the addresses to report listen - ([24a4923](https://github.com/Friends-Of-Noso/NosoGo/commit/24a49234bacd8cf8123414793425834f2793443a))
- Using `DatabasePath` on config - ([05aa952](https://github.com/Friends-Of-Noso/NosoGo/commit/05aa9525782094243f8924ff479811d07da6147c))
- Some light changes to version - ([5db787c](https://github.com/Friends-Of-Noso/NosoGo/commit/5db787c86c5761c0c9b6c0f9867d338df7782adc))
- No need for capacity field - ([f5fec22](https://github.com/Friends-Of-Noso/NosoGo/commit/f5fec22e83445eb16ccf81efe5bde461b41fcf93))

### 📚 Documentation

- *(README)* Correct branch for license link - ([f665d38](https://github.com/Friends-Of-Noso/NosoGo/commit/f665d38a42ca080de4eb7191c8a55805dde5c728))
- *(README)* Fix build badge - ([408a809](https://github.com/Friends-Of-Noso/NosoGo/commit/408a80916a92743fbf8a1a2652af7aeb9d5b931c))
- *(README)* Adding badges - ([bb7afd8](https://github.com/Friends-Of-Noso/NosoGo/commit/bb7afd8845cb530ff381373ebff64406fb37a5c3))
- *(README)* Match the repo description in the file - ([73edef8](https://github.com/Friends-Of-Noso/NosoGo/commit/73edef881e7ff904a0d0e7b644946cdee84284dd))

### ⚙️ Miscellaneous Tasks

- Adding the `MIT` license - ([9c51b01](https://github.com/Friends-Of-Noso/NosoGo/commit/9c51b01cb825a3989269828b3c6cf74f4405bdf2))
- `GOOS` for Windows - ([e766a4c](https://github.com/Friends-Of-Noso/NosoGo/commit/e766a4c524dbbd12b33aac5f8b40fbeb1631ad1e))
- Using `matrix` - ([8cfccb7](https://github.com/Friends-Of-Noso/NosoGo/commit/8cfccb715c87ea5b4de578fa8ffea51244f139e7))
- Adding step for `go mod tidy` - ([ae4e3eb](https://github.com/Friends-Of-Noso/NosoGo/commit/ae4e3eb1236b35bd461c9db01c4bc86d23e78778))
- Add Protobuf compiler - ([692c606](https://github.com/Friends-Of-Noso/NosoGo/commit/692c6067bebca149e61259f80daec9f71b8fb6ce))
- Adding `nosogocli` build - ([b8b446d](https://github.com/Friends-Of-Noso/NosoGo/commit/b8b446dbd6035dbddfd6a5d2e14a269e3d0be625))
- Removing erroneous created files - ([56a8526](https://github.com/Friends-Of-Noso/NosoGo/commit/56a8526a88ed98fa34bef1ff7dad22c5fc601959))
- Moving tests to single folder - ([fc28a20](https://github.com/Friends-Of-Noso/NosoGo/commit/fc28a2019eeeac9219772a303de35278a55b8a6b))
- Cleaning up `Makefile` - ([d9ced37](https://github.com/Friends-Of-Noso/NosoGo/commit/d9ced3793f4f01a7f5507b710c1979890b960746))
- Replacing some imports - ([c08b227](https://github.com/Friends-Of-Noso/NosoGo/commit/c08b22713211506c0489a6c89dbf6a96824afd2c))
- Renaming folder - ([05f432e](https://github.com/Friends-Of-Noso/NosoGo/commit/05f432e08cf6ba1b6eca01dfe443d2ef7f74fd12))
- Go modules init - ([86593b7](https://github.com/Friends-Of-Noso/NosoGo/commit/86593b72b49100f5a87d656c46ea9e2b5220c607))
- Initial commit - ([10c000c](https://github.com/Friends-Of-Noso/NosoGo/commit/10c000c5bd70ea10ae36cae795c3f52596e2b7f5))

## New Contributors ❤️

* [@gcarreno](https://github.com/gcarreno) made their first contribution

