Feature: Console

    Scenario: Text
        When I write to console:
        """
        Hello World!

        I'm John
        """

        Then console output is:
        """
        Hello World!

        I'm John
        """

    Scenario: JSON
        When I write to console:
        """
        [
            {
                "username": "user@example.org",
                "password": "123456"
            }
        ]
        """

        Then console output is:
        """
        [
            {
                "username": "user@example.org",
                "password": "123456"
            }
        ]
        """

    Scenario: Regex
        When I write to console:
        """
        [
            {
                "username": "user@example.org",
                "password": "123456"
            }
        ]
        """

        Then console output matches:
        """
        \[
            {
                "username": "user@example.org",
                "password": "[0-9]+"
            }
        \]
        """
