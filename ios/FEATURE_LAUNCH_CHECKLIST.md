# iOS App Feature Launch Checklist

This checklist validates that the iOS app implements all features described in the meal planning system documentation. Use this to ensure feature completeness before launch.

## Legend
- ✅ = Feature implemented and validated
- ⚠️ = Feature partially implemented or needs review
- ❌ = Feature not implemented
- N/A = Feature not applicable to iOS app (backend-only)

---

## 1. Meal Planning Core Features

### 1.1 Meal Plan Creation
- [x] **Create meal plan with name and time period**
  - [x] User can set meal plan name
  - [x] User can set voting deadline
  - [x] Meal plan can be created ad-hoc (at will)
  - [x] Validation: Name is required

- [x] **Add events to meal plan**
  - [x] User can add multiple events to a meal plan
  - [x] Each event has meal name (breakfast, lunch, dinner, snack)
  - [x] Each event has start and end times
  - [x] Each event can have notes
  - [x] Validation: At least one event required

- [ ] **Add meal options to events**
  - [x] User can search for existing meals to add as options
  - [x] Multiple options can be added per event
  - [x] Each option references a meal (collection of recipes)
  - [x] User can set meal scale for each option
  - [ ] ~~User can assign cook/dishwasher roles (optional)~~
  - [ ] ~~User can add notes to options (optional)~~

- [x] **Recipe option selections during creation**
  - [x] User can specify preferences for alternative ingredients/instruments/vessels
  - [x] Selections can be made for option groups in recipes
  - [x] User can skip selections (defaults to optionIndex: 0)
  - [x] Selections are saved with meal plan creation

### 1.2 Voting Phase
- [x] **View meal plans awaiting votes**
  - [x] Meal plans with pending votes are displayed
  - [x] Voting deadline is clearly shown
  - [x] Status indicates "awaiting votes"

- [x] **Rank meal plan options**
  - [x] User can rank options in order of preference for each event
  - [x] Drag-and-drop reordering is supported
  - [x] Visual feedback shows current ranking
  - [x] All options must be ranked before submission

- [x] **Change votes before deadline**
  - [x] User can modify rankings until deadline passes
  - [x] Previous votes are loaded and can be edited
  - [x] Changes are saved when submitted

- [x] **Abstain from voting**
  - [x] User can abstain from voting on some options while voting on others
  - [x] Abstention is clearly indicated in UI
  - [x] Abstention is properly submitted to backend

- [x] **Track voting status**
  - [x] System shows who has and hasn't voted
  - [x] Voting status is visible to all account members
  - [x] Visual indicators show voting completion status

- [x] **Lock ballots before submission**
  - [x] User can lock individual event ballots
  - [x] Locked ballots cannot be reordered
  - [x] All ballots must be locked before submission

- [x] **Submit votes**
  - [x] Submit button only appears when all ballots are complete and locked
  - [x] Votes are submitted with proper rank values (0 = first choice)
  - [x] Success/error feedback is provided
  - [x] User is redirected after successful submission

### 1.3 Finalization & Results
- [ ] **View finalized meal plans**
  - [ ] Finalized meal plans are displayed separately from pending
  - [ ] Finalization status is clearly indicated
  - [ ] Winning options are highlighted for each event
  - [ ] Tie-broken options are marked (if applicable)

- [ ] **View finalization results**
  - [ ] Winning meal option is shown for each event
  - [ ] Results are displayed after voting deadline passes
  - [ ] Results reflect Schulze voting method outcomes

### 1.4 Meal Plan Details
- [ ] **View meal plan details**
  - [ ] All events are displayed with their options
  - [ ] Event times and notes are shown
  - [ ] Meal plan name and voting deadline are displayed
  - [ ] Status (awaiting votes, finalized) is shown

- [ ] **Navigate to related features**
  - [ ] Link to grocery list from finalized meal plan
  - [ ] Link to task list from finalized meal plan
  - [ ] Link to vote interface from pending meal plan

---

## 2. Recipe Features

### 2.1 Recipe Discovery
- [ ] **View recipe list**
  - [ ] All available recipes are displayed
  - [ ] Recipes are filtered by status (e.g., "submitted")
  - [ ] Loading states are handled
  - [ ] Error states are handled
  - [ ] Empty states are handled

- [ ] **Recipe search** (if implemented)
  - [ ] User can search recipes by name
  - [ ] Search results are displayed
  - [ ] Search is debounced/throttled appropriately

- [ ] **View recipe details**
  - [ ] Recipe name, description, and source are displayed
  - [ ] Estimated portions (min/max) are shown
  - [ ] Portion names (singular/plural) are displayed
  - [ ] Component type (main, side, dessert, etc.) is shown
  - [ ] Created/updated timestamps are displayed

### 2.2 Recipe Steps & Instructions
- [ ] **View recipe steps**
  - [ ] Steps are displayed in order
  - [ ] Each step shows preparation method
  - [ ] Step notes/instructions are displayed
  - [ ] Step index/numbering is clear

- [ ] **View step ingredients**
  - [ ] All ingredients for each step are listed
  - [ ] Ingredient names are displayed
  - [ ] Quantities (min/max ranges) are shown
  - [ ] Measurement units are displayed
  - [ ] Option groups are indicated (if applicable)

- [ ] **View step instruments**
  - [ ] Required instruments are listed per step
  - [ ] Instrument names are displayed
  - [ ] Option groups for instruments are shown (if applicable)

- [ ] **View step vessels**
  - [ ] Required vessels are listed per step
  - [ ] Vessel names are displayed
  - [ ] Option groups for vessels are shown (if applicable)

- [ ] **View step products**
  - [ ] Products (outputs) of each step are displayed
  - [ ] Discrete products show item quantity and per-item measurement
  - [ ] Continuous products show total measurement quantity
  - [ ] Product types (ingredient, instrument, vessel) are indicated

### 2.3 Recipe Option Groups
- [ ] **Display option groups**
  - [ ] Alternative ingredients are grouped together
  - [ ] Alternative instruments are grouped together
  - [ ] Alternative vessels are grouped together
  - [ ] Primary option (optionIndex: 0) is clearly indicated
  - [ ] Alternatives are visually grouped

- [ ] **Select options during recipe viewing**
  - [ ] User can select which alternative to use
  - [ ] Selection affects displayed quantities
  - [ ] Default selection (optionIndex: 0) is applied if none chosen

### 2.4 Recipe Scaling
- [ ] **Scale recipe quantities**
  - [ ] User can adjust recipe scale (e.g., 2x, 0.5x)
  - [ ] Ingredient quantities are multiplied by scale factor
  - [ ] Discrete products: item quantity scales, per-item measurement stays constant
  - [ ] Continuous products: total measurement quantity scales
  - [ ] Scale factor is clearly displayed
  - [ ] Quantities update in real-time as scale changes

### 2.5 Recipe Relationships
- [ ] **View associated recipes**
  - [ ] Recipes used as ingredients are displayed
  - [ ] Associated recipes are shown in a flat list
  - [ ] User can navigate to associated recipe details
  - [ ] Associated recipes don't show nested associations (flattened)

- [ ] **Recipe cloning** (if implemented)
  - [ ] User can clone an existing recipe
  - [ ] Cloned recipe references original via InspiredByRecipeID
  - [ ] User becomes creator of cloned recipe
  - [ ] All IDs are regenerated appropriately

### 2.6 Recipe Prep Tasks
- [ ] **View prep tasks**
  - [ ] Prep tasks are displayed separately from cooking steps
  - [ ] Prep tasks show what can be done ahead of time
  - [ ] Prep tasks are linked to their recipes

### 2.7 Recipe Media
- [ ] **Display recipe media** (if implemented)
  - [ ] Recipe-level images/videos are displayed
  - [ ] Step-level images/videos are displayed
  - [ ] Media is properly loaded and cached

---

## 3. Meal Features

### 3.1 Meal Discovery
- [ ] **Search for meals**
  - [ ] User can search meals by name
  - [ ] User can filter by component type
  - [ ] Search results are displayed
  - [ ] Search is used when creating meal plan options

- [ ] **View meal details**
  - [ ] Meal name and description are displayed
  - [ ] Estimated portions (min/max) are shown
  - [ ] All components are listed with their types
  - [ ] Component types are clearly labeled (main, side, appetizer, etc.)
  - [ ] Recipe scales for each component are shown
  - [ ] User can navigate to component recipe details

### 3.2 Meal Components
- [ ] **View meal components**
  - [ ] All components are displayed in order
  - [ ] Component types are shown (main, side, dessert, etc.)
  - [ ] At least one "main" component is present
  - [ ] Recipe scale for each component is displayed
  - [ ] User can view individual recipe details from components

### 3.3 Meal Scaling
- [ ] **Scale meals in meal plans**
  - [ ] User can set meal scale when adding to meal plan
  - [ ] Meal scale affects all component recipes proportionally
  - [ ] Scale factor is clearly displayed
  - [ ] Quantities update based on meal scale

---

## 4. Execution Features (Post-Finalization)

### 4.1 Grocery Lists
- [ ] **View grocery lists**
  - [ ] Grocery lists are available for finalized meal plans
  - [ ] Items are grouped by status (needs, acquired, already owned, unavailable)
  - [ ] Items show ingredient name, quantity needed, and measurement unit
  - [ ] Items reflect recipe option selections (only selected alternatives appear)
  - [ ] Items reflect meal and recipe scaling

- [ ] **Modify grocery list items**
  - [ ] User can mark items as acquired
  - [ ] User can mark items as already owned
  - [ ] User can mark items as needs
  - [ ] User can mark items as unavailable
  - [ ] User can edit quantity needed (min/max)
  - [ ] User can update quantity purchased
  - [ ] Changes are saved to backend
  - [ ] UI updates reflect changes immediately

- [ ] **Grocery list organization**
  - [ ] Items can be filtered by status
  - [ ] Items are sorted appropriately
  - [ ] Similar ingredients are consolidated (if implemented)

### 4.2 Task Management
- [ ] **View tasks**
  - [ ] Tasks are displayed for finalized meal plans
  - [ ] Tasks are grouped by type (prep, cooking)
  - [ ] Tasks show assigned user (if applicable)
  - [ ] Tasks show completion status
  - [ ] Tasks are linked to their recipes

- [ ] **Task details**
  - [ ] Task shows recipe name and step information
  - [ ] Task shows due date/time
  - [ ] Task shows assigned user
  - [ ] Task shows completion status

- [ ] **Complete tasks**
  - [ ] User can mark tasks as complete
  - [ ] User can mark individual steps as complete (for multi-step tasks)
  - [ ] Completion status is saved to backend
  - [ ] UI updates reflect completion

- [ ] **Task subtasks**
  - [ ] Multi-step tasks show subtasks (recipe steps)
  - [ ] User can expand/collapse task details
  - [ ] User can mark individual subtasks as complete
  - [ ] Parent task completion reflects subtask completion

### 4.3 Recipe Performance
- [ ] **View recipe for cooking**
  - [ ] Recipe steps are displayed in order
  - [ ] Current step is highlighted
  - [ ] User can navigate between steps
  - [ ] Ingredients, instruments, and vessels are shown per step
  - [ ] Quantities reflect scaling and selections

- [ ] **Track cooking progress**
  - [ ] User can mark steps as complete
  - [ ] Progress is saved
  - [ ] Completed steps are visually indicated

---

## 5. User Experience & Edge Cases

### 5.1 Error Handling
- [ ] **Network errors**
  - [ ] Network failures are caught and displayed
  - [ ] User can retry failed operations
  - [ ] Error messages are user-friendly

- [ ] **Validation errors**
  - [ ] Form validation errors are displayed
  - [ ] Field-level errors are shown
  - [ ] User can correct errors and resubmit

- [ ] **Backend errors**
  - [ ] API errors are handled gracefully
  - [ ] Error messages are displayed to user
  - [ ] User can recover from errors

### 5.2 Loading States
- [ ] **Loading indicators**
  - [ ] Loading states are shown during API calls
  - [ ] Progress indicators are appropriate
  - [ ] Users can't interact with loading content

### 5.3 Data Refresh
- [ ] **Refresh data**
  - [ ] Pull-to-refresh is supported where appropriate
  - [ ] Data refreshes after mutations (create, update, delete)
  - [ ] Data refreshes when returning to screens

### 5.4 Navigation
- [ ] **Navigation flow**
  - [ ] User can navigate between all major screens
  - [ ] Back navigation works correctly
  - [ ] Deep linking works (if implemented)
  - [ ] Navigation state is preserved appropriately

### 5.5 Edge Cases
- [ ] **Empty states**
  - [ ] Empty meal plan list is handled
  - [ ] Empty recipe list is handled
  - [ ] Empty grocery list is handled
  - [ ] Empty task list is handled
  - [ ] Empty search results are handled

- [ ] **Deadline handling**
  - [ ] Voting deadline passed is handled
  - [ ] Meal plan finalization is reflected in UI
  - [ ] Users can't vote after deadline (if enforced)

- [ ] **Tie-breaking**
  - [ ] Tie-broken options are marked appropriately
  - [ ] UI indicates when a tie was broken

- [ ] **No votes scenario**
  - [ ] When no one votes, system handles appropriately
  - [ ] One option is chosen and marked as tie-broken

---

## 6. Data Integrity & Validation

### 6.1 Meal Plan Validation
- [ ] **Required fields**
  - [ ] Meal plan name is required
  - [ ] Voting deadline is required
  - [ ] At least one event is required
  - [ ] Each event has at least one option

- [ ] **Data consistency**
  - [ ] Meal references are valid
  - [ ] Recipe references in meals are valid
  - [ ] Option selections reference valid recipes/steps

### 6.2 Recipe Validation
- [ ] **Recipe structure**
  - [ ] Recipes have at least 2 steps
  - [ ] Each step has at least one instrument OR vessel
  - [ ] Required fields are present (name, portions, etc.)

- [ ] **Option groups**
  - [ ] Option groups are properly structured
  - [ ] Selected options are valid

### 6.3 Meal Validation
- [ ] **Meal structure**
  - [ ] At least one component has type "main"
  - [ ] All components have valid component types
  - [ ] Recipe references are valid

---

## 7. Performance & Polish

### 7.1 Performance
- [ ] **Loading performance**
  - [ ] Lists load efficiently
  - [ ] Images/media load efficiently
  - [ ] Search is responsive

- [ ] **Memory management**
  - [ ] No memory leaks
  - [ ] Large lists are handled efficiently (lazy loading)

### 7.2 UI/UX Polish
- [ ] **Visual design**
  - [ ] Consistent styling throughout app
  - [ ] Colors and typography are consistent
  - [ ] Icons are appropriate and consistent

- [ ] **Accessibility**
  - [ ] VoiceOver support (if applicable)
  - [ ] Dynamic type support (if applicable)
  - [ ] Color contrast meets standards

- [ ] **Animations**
  - [ ] Transitions are smooth
  - [ ] Loading animations are appropriate
  - [ ] State changes are animated where appropriate

---

## 8. Documentation-Specific Features

### 8.1 Features Marked as "Unimplemented" in Docs
- [ ] **EligibleForMeals flag** (docs say largely unimplemented)
  - [ ] If implemented: UI respects this flag when creating meals
  - [ ] If not implemented: Document that it's not enforced

- [ ] **EligibleForMealPlans flag** (docs say largely unimplemented)
  - [ ] If implemented: UI respects this flag when creating meal plans
  - [ ] If not implemented: Document that it's not enforced

- [ ] **Slug-based routing** (docs say not implemented)
  - [ ] If implemented: Recipes can be accessed via slug
  - [ ] If not implemented: Use recipe ID for routing

### 8.2 Future Features (Not Required for Launch)
- [ ] **Instant runoff voting** (docs say TODO)
  - [ ] Not required for initial launch
  - [ ] Document as future enhancement

- [ ] **Recipe filtering by ingredient preferences** (docs say TODO)
  - [ ] Not required for initial launch
  - [ ] Document as future enhancement

- [ ] **Vector search for meals** (docs mention future plans)
  - [ ] Not required for initial launch
  - [ ] Current search implementation is sufficient

---

## Testing Checklist

### 8.1 Manual Testing Scenarios
- [ ] **Complete meal plan flow**
  1. Create meal plan with multiple events
  2. Add meal options to each event
  3. Set recipe option selections
  4. Submit meal plan
  5. Vote on meal plan
  6. Wait for finalization (or manually finalize)
  7. View grocery list
  8. View tasks
  9. Complete tasks
  10. Mark grocery items as acquired

- [ ] **Recipe viewing flow**
  1. Browse recipe list
  2. View recipe details
  3. Scale recipe
  4. Select option group alternatives
  5. View associated recipes
  6. Navigate to associated recipe

- [ ] **Meal search and selection**
  1. Search for meals
  2. View meal details
  3. Add meal to meal plan option
  4. Set meal scale

- [ ] **Error scenarios**
  1. Network failure during API call
  2. Invalid form submission
  3. Voting after deadline (if enforced)
  4. Accessing deleted/non-existent resources

### 8.2 Edge Case Testing
- [ ] Meal plan with single event
- [ ] Meal plan with many events
- [ ] Event with single option
- [ ] Event with many options
- [ ] Recipe with no option groups
- [ ] Recipe with many option groups
- [ ] Meal with single component
- [ ] Meal with many components
- [ ] Grocery list with no items
- [ ] Grocery list with many items
- [ ] Task list with no tasks
- [ ] Task list with many tasks

---

## Notes

- This checklist is based on the documentation in `docs/meal_planning.md`, `docs/meals.md`, and `docs/recipes.md`
- Features marked as "TODO" or "largely unimplemented" in the docs may not be required for launch
- Some features may be backend-only and not applicable to iOS app validation
- Update this checklist as features are implemented or requirements change
