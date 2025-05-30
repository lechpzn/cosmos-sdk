<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section from newest to oldest.
Each entry must include the Github issue reference in the following format:

* [#<issue-number>] Changelog message.

-->

# Changelog

## [Unreleased]

## [v1.0.0](https://github.com/cosmos/cosmos-sdk/releases/tag/api/v1.0.0)

### Breaking Changes

- The `api` module now depends on the `cometbft` buf instead of re-exporting it or the legacy tendermint protobuf APIs

## [v0.9.0](https://github.com/cosmos/cosmos-sdk/releases/tag/api/v0.9.0) - 2025-03-31

### Features

* [#23933](https://github.com/cosmos/cosmos-sdk/pull/23933) `x/protocolpool` API files
* [#23815](https://github.com/cosmos/cosmos-sdk/pull/23815) `x/epochs` API files
* [#23708](https://github.com/cosmos/cosmos-sdk/pull/23708) `unordered` transaction support

### Improvements

* [#24227](https://github.com/cosmos/cosmos-sdk/pull/24227) Minor dependency bumps


