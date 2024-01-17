import os
import re
import sys
from typing import Dict, List

from entry import Entry
from change_type import ChangeType
from release import Release

# Allowed release pattern: vX.Y.Z(-rcN) (YYYY-MM-DD)
RELEASE_PATTERN = re.compile(
    r'^## (Unreleased|\[(?P<version>v\d+\.\d+\.\d+(-rc\d+)?)] - \d{4}-\d{2}-\d{2})$',
)


class Changelog:
    """
    This class represents the contents of the changelog and provides methods to parse it.
    """

    def __init__(self, filename: str):
        self.contents: List[str]
        self.filename: str = filename

        self.problems: List[str] = []
        # TODO: extract releases type
        self.releases: Dict[str, Dict[str, Dict[int, Dict[str, str]]]] = {}

        if not os.path.exists(self.filename):
            raise FileNotFoundError(f'Changelog file "{self.filename}" not found')

        with open(self.filename, 'r') as file:
            self.contents = file.read()

    def parse(self) -> bool:
        """
        This function parses the changelog and checks if the structure is as expected.
        """

        current_release = None
        current_category = None

        for line in self.contents.split('\n'):
            # Check for Header 2 (##) to identify releases
            stripped_line = line.strip()
            if stripped_line[:3] == '## ':
                release = Release(line)
                release.parse()
                current_release = release.version
                if current_release in self.releases:
                    self.problems.append(f'Release "{current_release}" is duplicated in the changelog')
                else:
                    self.releases[current_release] = {}
                self.problems.extend(release.problems)
                continue

            # Check for Header 3 (###) to identify change types
            if stripped_line[:4] == '### ':
                change_type = ChangeType(line)
                change_type.parse()
                current_category = change_type.type
                if current_category in self.releases[current_release]:
                    self.problems.append(f'Change type "{current_category}" is duplicated in {current_release}')
                else:
                    self.releases[current_release][current_category] = {}
                self.problems.extend(change_type.problems)
                continue

            # Check for individual entries
            if stripped_line[:2] != '- ':
                continue

            # TODO: order by extending the types by entries and then process afterwards within each release to have sorted output.
            entry = Entry(line)
            entry.parse()
            self.problems.extend(entry.problems)

            self.releases[current_release][current_category][entry.pr_number] = {
                "description": entry.description
            }

        return self.problems == []


if __name__ == "__main__":
    changelog = Changelog(sys.argv[1])
    passed = changelog.parse()
    if not passed:
        print(f'Changelog file is not valid - check the following {len(changelog.problems)} problems:\n')
        print('\n'.join(changelog.problems))
        sys.exit(1)
