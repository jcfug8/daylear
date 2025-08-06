# Visibility
VisibilityLevel:
- VISIBILITY_LEVEL_UNSPECIFIED - ?? Is this the same as public?
- VISIBILITY_LEVEL_PUBLIC - The recipe is visible in the public recipe list.
- VISIBILITY_LEVEL_RESTRICTED - The recipe is visible only users that have: explicit access to the recipe or access to a circle that has access to the recipe.
- VISIBILITY_LEVEL_PRIVATE - The recipe is only visible to users that have explicit access to the recipe.
- VISIBILITY_LEVEL_HIDDEN - The recipe is only visible to the creator of the recipe.

## Standard Access:
This comes into play when you have explicit access to a resource. It directly follows the visibility level of the resource + your access.

## Delegated Access:
When viewing a user/circle's resources, this comes into play when you are viewing resources that you don't have explicit access to, but you have explicit access to their parent. It directly follows the visibility level of the resource + the parent's access.


### Get (Unified)
#### Process
- Get Resource
- if public
    - return resource
- if restricted
    - must have standard access or delegated user/circle access
- if private
    - must have standard access or delegated circle access
- if hidden
    - must have standard access

### List
#### What access to return
Will need:
- List (for authUser)
- ListDelegated (for parent)
    - May split this into two: one for circle and one for user
#### Auth User
##### Domain
##### Database

#### Circle Parent
##### Domain
##### Database

#### User Parent
##### Domain
##### Database