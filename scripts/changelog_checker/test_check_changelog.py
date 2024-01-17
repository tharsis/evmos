import os
import pytest
from shutil import copyfile

from check_changelog import Changelog, Entry

# Get the directory of this script
SCRIPT_DIR = os.path.dirname(os.path.realpath(__file__))


@pytest.fixture
def create_tmp_copy():
    tmp_file = os.path.join(SCRIPT_DIR, "testdata", "changelog_tmp.md")
    copyfile(
        os.path.join(SCRIPT_DIR, "testdata", "changelog_fail.md"),
        tmp_file,
    )
    yield tmp_file
    os.remove(tmp_file)


class TestParseChangelog:
    """
    This class collects all tests that are actually parsing dummy changelogs stored in
    markdown files in the testdata directory.
    """

    def test_pass(self):
        expected_result = {
            'Unreleased': {
                'State Machine Breaking': {
                    1922: {'description': 'Add `secp256r1` curve precompile.'},
                    1949: {'description': 'Add `ClaimRewards` custom transaction.'},
                },
                'API Breaking': {
                    2015: {'description': 'Rename `inflation` module to `inflation/v1`.'},
                    2078: {'description': 'Deprecate legacy EIP-712 ante handler.'},
                },
                'Improvements': {
                    1864: {
                        'description':
                            'Add `--base-fee` and `--min-gas-price` flags.',
                    },
                    1912: {'description': 'Add Stride outpost interface and ABI.'},
                },
                'Bug Fixes': {
                    1801: {'description': 'Fixed the problem `gas_used` is 0.'},
                },
            },
            'v15.0.0': {
                'API Breaking': {
                    1862: {'description': 'Add Authorization Grants to the Vesting extension.'},
                },
            },
        }

        changelog = Changelog(os.path.join(SCRIPT_DIR, "testdata", "changelog_ok.md"))
        assert changelog.parse() is True
        assert changelog.problems == [], "expected no failed entries"
        assert changelog.releases == expected_result, "expected different parsed result"

    def test_fail(self):
        changelog = Changelog(os.path.join(SCRIPT_DIR, "testdata", "changelog_fail.md"))
        assert changelog.parse() is False
        assert changelog.problems == [
            'PR link is not matching PR number 1948: "https://github.com/evmos/evmos/pull/1949"',
            '"ABI" should be used instead of "ABi"',
            '"outpost" should be used instead of "Outpost"',
            'PR description should end with a dot: "Fixed the problem `gas_used` is 0"',
            '"Invalid Category" is not a valid change type',
            'Change type "Bug Fixes" is duplicated in Unreleased',
            'Release "v15.0.0" is duplicated in the changelog',
            'Change type "API Breaking" is duplicated in v15.0.0',
        ]

    def test_fix(self, create_tmp_copy):
        changelog = Changelog(create_tmp_copy)
        assert changelog.parse(fix=True) is False
        assert changelog.problems == [
            'PR link is not matching PR number 1948: "https://github.com/evmos/evmos/pull/1949"',
            '"ABI" should be used instead of "ABi"',
            '"outpost" should be used instead of "Outpost"',
            'PR description should end with a dot: "Fixed the problem `gas_used` is 0"',
            '"Invalid Category" is not a valid change type',
            'Change type "Bug Fixes" is duplicated in Unreleased',
            'Release "v15.0.0" is duplicated in the changelog',
            'Change type "API Breaking" is duplicated in v15.0.0',
        ]

        # Here we parse the fixed changelog again and check that the automatic fixes were applied.
        fixed_changelog = Changelog(changelog.filename)
        assert fixed_changelog.parse(fix=False) is False
        assert fixed_changelog.problems == [
            '"Invalid Category" is not a valid change type',
            'Change type "Bug Fixes" is duplicated in Unreleased',
            'Release "v15.0.0" is duplicated in the changelog',
            'Change type "API Breaking" is duplicated in v15.0.0',
        ]

    def test_parse_changelog_nonexistent_file(self):
        with pytest.raises(FileNotFoundError):
            Changelog(os.path.join(SCRIPT_DIR, "testdata", "nonexistent_file.md"))
