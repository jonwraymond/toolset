# Changelog

## [1.0.1](https://github.com/jonwraymond/toolset/compare/v1.0.0...v1.0.1) (2026-01-29)

### Chores

- bump tooladapter dependency to v0.2.0

## [1.0.0](https://github.com/jonwraymond/toolset/compare/toolset-v0.2.1...toolset-v1.0.0) (2026-01-29)


### ⚠ BREAKING CHANGES

* **toolset:** ExportWithWarnings signature changed from   ([]any, []FeatureLossWarning) to ([]any, []FeatureLossWarning, []error)

### Features

* **toolset:** add Builder pattern ([f30d75b](https://github.com/jonwraymond/toolset/commit/f30d75b3b9bbaca745a86dd89c184151a1c527f3))
* **toolset:** add Exposure helpers ([0aa8ca5](https://github.com/jonwraymond/toolset/commit/0aa8ca5e374f29c437410f0a74c74d1b0c11fce1))
* **toolset:** add filter helpers ([c4abb05](https://github.com/jonwraymond/toolset/commit/c4abb052e8f6c49f08a6231d17988864e16d439c))
* **toolset:** add Policy interface ([661ede5](https://github.com/jonwraymond/toolset/commit/661ede55ce5b589b283a30f09565f61c0f586ed3))
* **toolset:** add Toolset core and package documentation ([c1fb410](https://github.com/jonwraymond/toolset/commit/c1fb4106aaad6965e89af67f60aea9358084fa5a))


### Bug Fixes

* **toolset:** remove local replace ([96830bc](https://github.com/jonwraymond/toolset/commit/96830bc9696de7566c85582101b697dfb0b2dd88))
* **toolset:** surface conversion errors ([bb4375f](https://github.com/jonwraymond/toolset/commit/bb4375fb90b35a867cacf07485022e59bb91bbf2))
* **toolset:** surface conversion errors in ExportWithWarnings ([8620257](https://github.com/jonwraymond/toolset/commit/8620257327b4b58bcb0fec9c094b268fd3479b76))

## [0.2.1](https://github.com/jonwraymond/toolset/compare/v0.2.0...v0.2.1) (2026-01-29)

### Bug Fixes

- Remove local replace directive to unblock CI module downloads.

## [0.2.0](https://github.com/jonwraymond/toolset/compare/v0.1.0...v0.2.0) (2026-01-29)

### Breaking Changes

- `ExportWithWarnings` now returns a third value (`[]error`) for conversion failures.

### Bug Fixes

- Surface conversion failures instead of silently dropping tools.
- Align feature loss detection with tooladapter recursion semantics.

## [0.1.0](https://github.com/jonwraymond/toolset/releases/tag/v0.1.0) (2026-01-29)

### Features

- add initial toolset docs and mkdocs configuration
- add initial module scaffold and version matrix
