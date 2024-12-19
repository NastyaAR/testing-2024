Feature: Two-factor authentification

  Scenario: Successful authentification
    Given I register with @N_AR24, password, client
    When I log with user_id and password
      Then I get message 'we send code to your telegram'
    When I provide user_id and code
      Then I get token 

  Scenario: Authentification failure (no such user)
    When I log without registration with user_id and password
      Then I get error 

  Scenario: Authentification failure (not valid verification code)
    When I login with user_id and password
      Then I get message 'we send code to your telegram'
    When I incorrectly provide user_id and code
      Then I get error 
  
  Scenario: Successful changing password
    Given I register with @N_AR24, password, client
    When I log with user_id and password
      Then I get message 'we send code to your telegram'
    When I want change to newpassword, user_id, code

  Scenario: Failure changing password
    Given I register with @N_AR24, password, client
    When I log with user_id and password
      Then I get message 'we send code to your telegram'
    When I want change to empty password, user_id, code

