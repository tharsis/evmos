import os
import pytest
from check_changelog import parse_changelog

# Get the directory of this script
SCRIPT_DIR = os.path.dirname(os.path.realpath(__file__))


def test_parse_changelog_ok():
    expected_result = {
        'Unreleased': {
            'State Machine Breaking': {
                1922: {'description': 'Add `secp256r1` curve precompile.'},
                1949: {'description': 'Add `ClaimRewards` custom transaction.'},
            },
            'API Breaking': {
                2015: {'description': 'Rename `inflation` module to `inflation/v1`.'},
                2078: {'description': 'Deprecate legacy EIP712 ante handler.'},
            },
            'Improvements': {
                1864: {
                    'description':
                        'Add `--base-fee` and `--min-gas-price` flags to the command `evmosd testnet init-files`.',
                },
                1912: {'description': 'Add Stride Outpost interface and ABI.'},
            },
            'Bug Fixes': {
                1801: {'description': 'Fixed the problem gas_used is 0.'},
            },
        },
        'v15.0.0': {
            'API Breaking': {
                1862: {'description': 'Add Authorization Grants to the Vesting extension.'},
            },
        },
    }

    releases, failed_entries = parse_changelog(os.path.join(SCRIPT_DIR, "testdata", "changelog_ok.md"))
    assert failed_entries == []
    assert releases == expected_result


def test_parse_changelog_invalid_pr_link():
    expected_result = {
        'Unreleased': {
            'State Machine Breaking': {
                1948: {'description': 'Add `ClaimRewards` custom transaction.'},
            },
        },
    }
    releases, failed_entries = parse_changelog(
        os.path.join(SCRIPT_DIR, "testdata", "changelog_invalid_entry_pr_not_in_link.md"),
    )
    assert failed_entries == [
        'Invalid PR link in Unreleased - State Machine Breaking - 1948: "- (distribution-precompile) [#1948]('
        'https://github.com/evmos/evmos/pull/1949) Add `ClaimRewards` custom transaction."',
    ]
    assert releases == expected_result


def test_parse_changelog_malformed_description():
    expected_result = {
        'Unreleased': {
            'State Machine Breaking': {
                1949: {'description': 'add `ClaimRewards` custom transaction'},
            },
        },
    }
    releases, failed_entries = parse_changelog(
        os.path.join(SCRIPT_DIR, "testdata", "changelog_invalid_entry_misformatted_description.md"),
    )
    assert failed_entries == [
        'Invalid PR description in Unreleased - State Machine Breaking - 1949: "- (distribution-precompile) [#1949]('
        'https://github.com/evmos/evmos/pull/1949) add `ClaimRewards` custom transaction"',
    ]
    assert releases == expected_result


def test_parse_changelog_invalid_category():
    expected_result = {
        'Unreleased': {
            'Invalid Category': {
                1949: {'description': 'Add `ClaimRewards` custom transaction.'},
            },
        },
    }
    releases, failed_entries = parse_changelog(os.path.join(SCRIPT_DIR, "testdata", "changelog_invalid_category.md"))
    assert failed_entries == ["Invalid change category in Unreleased: \"Invalid Category\""]
    assert releases == expected_result


def test_parse_changelog_invalid_header():
    with pytest.raises(ValueError):
        parse_changelog(os.path.join(SCRIPT_DIR, "testdata", "changelog_invalid_version.md"))


def test_parse_changelog_invalid_date():
    with pytest.raises(ValueError):
        parse_changelog(os.path.join(SCRIPT_DIR, "testdata", "changelog_invalid_date.md"))


def test_parse_changelog_nonexistent_file():
    with pytest.raises(FileNotFoundError):
        parse_changelog("nonexistent_file.md")
