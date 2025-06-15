# Access
The main fields are:
- the "target" - the resource that is being accessed
- the "issuer" - user/circle that gave the access
- the "recipient" - user/circle that received the access
- the "status" - the status of the access
- the "access level" - the access level of the access
## Access Levels
- Unspecified - the recipient can view the target but only because the target is public.
- Read - the recipient can view the target and has explicit access to the target
- Write - the recipient can view and edit the target
- Admin - the recipient can view, edit, and delete the target
## Access Status
- Pending - waiting for the user to accept or deny the invitation (this override the access level and makes the recipient have read only until accepted)
- Accepted - the user has accepted the invitation



# Recipes
The Recipe resource should return all the recipe information plus the access data of the current account. This resource should be returned if you have explicit access to the recipe or if the recipe is public. If you only want to get recipes you have explicit access to, you should use the filter on the "access_level" field. This field will be used to allow the UI to determine what to show where (if different tabs are used for different access levels). It will also keep the recipe view not cluttered with all the public recipes.

The UI should show a UI to list recipes. I would expect there to be a table for public recipes, a table for their recipes, and a table for recipes the user has access to.
## Levels
- Unspecified
  - the recipe details can be viewed
- Read
  - the recipe details can be viewed
- Write
  - the recipe details can be viewed
  - the user can edit it
  - the user can manage access to the recipe
- Admin
  - the recipe details can be viewed
  - the user can edit it
  - the user can manage access to the recipe
  - the user can delete the recipe
### Persistence
- Unspecified - this is saved base on the public flag of the recipe
- All Other - Explicitly stored on some "access" table
### Recipe Access List
This resource contains a list of the users and circles that have access to the recipe separated by status. It should expose the ability to get the list, and add/remove access.

#### `/recipes/{recipe}/recipeAccessList` 
This holds lists of those that have access to a recipe.

This use case for this is when viewing a recipe and wanting to manage who has access to the recipe. This is shown on the UI as a modal that allows the user to manage the access.
- Get `/recipes/{recipe}/recipeAccessList`
  - Allows the user to get the list of access for a recipe.
- Add `/recipes/{recipe}/recipeAccessList`
  - Allows the user to add access to a recipe.
- Remove `/recipes/{recipe}/recipeAccessList`
  - Allows the user to remove access from a recipe.
- Accept - `/recipes/{recipe}/recipeAccessList`
  - Allows the user to accept an invitation to access a recipe.

#### `/recipeAccessList`
This holds lists of accesses to recipes that the current user/circle manages access for. (If you are actively acting as a circle, you can only manage the access for the circle and the same visa-versa when you are actively acting as a user you can only manage the access for the user. This resource will not cross over between user and circle.)

This use case for when a user wants to manage all of the accesses for recipes they manage or recipes they want access to. This is shown on the UI as a view that allows the user to manage the access. It would probably have two different tabs for recipes they're sharing, recipes they're requesting access to.
- Get `/recipeAccessList`
  - Allows the user to get the list of accesses for recipes they manage or recipes they want access to.
- Remove - `/recipeAccessList`
  - Allows the user to remove access from a recipe.
- Accept - `/recipeAccessList`
  - Allows the user to accept an invitation to access a recipe.


# Circles
The Circle resource should return the circle information plus the access data of the current account. This resource should be returned if you have access to the circle or if the circle is public. If you only want to get circles you have explicit access to, you should use the filter on the "access_level" field. This field will be used to allow the UI to determine what to show where (if different tabs are used for different access levels). It will also keep the circle view not cluttered with all the public circles.

THe UI should show a UI to list circles. I would expect there to be a table for public circles and a table for circles the user has access to.

I need to remove the current publicCircle resource because it is not needed. If a circle is public then you don't need access to view it. You can, however, ask for write access to a public circle.
## Levels
- Unspecified
  - the circle basic details can be viewed
- Read
  - the circle basic details can be viewed including the recipes
- Write
  - the circle basic details can be viewed including the recipes
  - the user can edit it
  - the user can manage access to the circle
- Admin
  - the circle basic details can be viewed including the recipes
  - the user can edit it
  - the user can manage access to the circle
  - the user can delete the circle
### Persistence
- Unspecified - this is saved base on the public flag of the circle
- All Other - Explicitly stored on some "access" table
## Circle Access List
This resource contains lists of the users that have access to the circle separated by status. It should expose the ability to get the list, and add/remove access.

### /circleAccessList
A list of access for the current account (circle). This is shown on the UI as a view that allows the user to manage the access. It would probably have two different tabs for users they're sharing, users they're requesting access to.

The use case for this is as a circle you want to manage what users have access to you. This will allow a user to manage who can access the circle.
- Get - show all accesses for the current account (circle)
  - /circleAccessList
- Add - Invite a user to access the circle (and ask for write access to the circle if it is public)
  - /circleAccessList
- Accept - Accept an invitation to access the circle
  - /circleAccessList
- Remove - Remove or deny access to the circle
  - /circleAccessList
####
### How does a user manage requesting to access, accepting access, or removing access a circle?
- How the users see it
  - To remove their access
    - they go to the circle settings and click a button to remove their access
    - they go to the circle list page, click on some access tab that should show their access list and they can remove their access
  - To accept access
    - they go to the circle list page, click on some access tab that should show their access list and they can accept the access
- The resource used
  - you have both circle and public circle resource. They should both return the access data of the current account. We can use that to determine what to show where


# Users
The User resource should return the access level the current account has for the user. This resource should be returned if you have access to the user or if the user is public. If you only want to get users you have explicit access to, you should use the filter on the "access_level" field. This field will will be used to allow the UI to determine what to show where (if different tabs are used for different access levels). It will also keep the circle view not cluttered with all the public users.
- Unspecified
  - the basic user details can be viewed
- Read
  - the basic user details can be viewed
- Write
  - they can share recipes with the user
- Admin
  - this is the user that is logged in
## User Access List
This resource contains lists of the users that have access to the user separated by status.
### users/{user}/accessList
- Get
  - users/{user}/access
- List
  - users/-/access