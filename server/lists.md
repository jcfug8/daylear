# Lists
## Functionality
- CRUD Lists
- CRUD ListItems
- List Sections
- Allow items to be completed
- Setting to show completed items
- Items can become auto uncompleted based on a RRULE
- Assign points to an item

## UI
- Make a list
- Add a list item
- move list items
- create a list section
- move a list sections
- Add a RRULE to a list item
- User a complete/uncomplete list item
- Show who completed an item
- Show completion history?
- Add a due date
- show when a due date has passed
- Add points to an item

## DB
- List
    - Id
    - Title
    - Description
    - Create Time
    - Update Time
    - Show Completed
    - Visibility
    - Sections map[int]string
- List Item
    - Id
    - List Id
    - Title
    - Create Time
    - Update Time
    - Points
    - List Section
    - Recurrence Rule
- List Item Completion
    - Id
    - List Item Id
    - User Id
    - List Id
    - Create Time
- ListAccess
    - Id
    - List Id
    - Requester User Id
    - Requester Circle Id
    - Recipient User Id
    - Recipient Circle Id
    - Permission Level
    - State
    - Create Time
    - Update Time

Point Sets