# AT Protocol Integration Guide

This document outlines how to integrate the AT Protocol (ATProto) into the Dinner Done Better meal planning application. The AT Protocol is the decentralized social networking protocol used by Bluesky and other federated social platforms.

## Table of Contents

1. [What This Means for Users](#what-this-means-for-users) ⭐ **START HERE**
2. [User Experience Examples](#user-experience-examples)
3. [Why Integrate AT Protocol?](#why-integrate-at-protocol)
4. [Integration Approaches](#integration-approaches)
5. [Architecture Design](#architecture-design)
6. [Implementation Steps](#implementation-steps)
7. [Go Libraries and Tools](#go-libraries-and-tools)
8. [Domain Models](#domain-models)
9. [API Design](#api-design)
10. [Security Considerations](#security-considerations)

## Do You Actually Need AT Protocol? (The Honest Answer)

### The Simple Alternative

You're absolutely right to question this. You could achieve 80% of the value with 20% of the effort:

**Simple Approach (No AT Protocol):**
1. Format a nice post with meal plan details
2. Generate a shareable link (e.g., `dinnerdonebetter.com/meal-plans/abc123`)
3. User copies the formatted text and link
4. User pastes it into Bluesky manually
5. Done.

**What you'd build:**
- A "Share" button that formats text and copies to clipboard
- A public link view for meal plans/recipes
- Maybe a QR code or deep link

**Development time:** 1-2 days  
**Complexity:** Minimal  
**User friction:** Medium (they have to copy/paste)

### When AT Protocol Integration Makes Sense

AT Protocol integration is **only worth it** if you want these specific features:

#### 1. **Automatic Posting (Reduces Friction)**
- **Without AT Protocol:** User clicks share → copies text → opens Bluesky → pastes → posts
- **With AT Protocol:** User clicks share → done (post appears automatically)

**Is this worth it?** Maybe. Studies show reducing friction increases sharing by 2-3x. But if your users are already motivated to share, the friction might not matter.

#### 2. **Discovery (Reading from Bluesky)**
- **Without AT Protocol:** Users only see content in your app
- **With AT Protocol:** Users can discover recipes shared on Bluesky, see a feed of meal plans from people they follow

**Is this worth it?** This is where AT Protocol shines. If you want users to discover content from the Bluesky network, you need AT Protocol to read that data.

#### 3. **Rich Embeds (Link Previews)**
- **Without AT Protocol:** Links show as plain URLs
- **With AT Protocol:** Links can show rich previews with images, structured data

**Is this worth it?** Probably not worth it alone, but nice to have.

#### 4. **Social Graph Integration**
- **Without AT Protocol:** Users can't follow each other's meal plans in-app
- **With AT Protocol:** Users can follow Bluesky users and see their meal plans in your app

**Is this worth it?** Only if you want the app to feel like a social network, not just a tool.

### The Real Question: What Do You Actually Want?

**If you just want users to share meal plans:**
→ **Don't use AT Protocol.** Build a simple share button that formats text and copies to clipboard. Add a public link view. Done in 1-2 days.

**If you want users to discover recipes from Bluesky:**
→ **You need AT Protocol** (or at least the read APIs). You can't discover Bluesky content without it.

**If you want automatic posting (one-click share):**
→ **You need AT Protocol** for the OAuth and posting APIs. But ask yourself: is the friction reduction worth 6-8 weeks of development?

**If you want a social feed in your app:**
→ **You need AT Protocol** to read the social graph and content.

### My Honest Recommendation

**Start with the simple approach:**
1. Build a share button that formats text
2. Add public link views for meal plans/recipes
3. See if users actually share
4. Measure engagement

**Then decide:**
- If users are sharing a lot → Consider AT Protocol for automatic posting
- If you want discovery features → You'll need AT Protocol
- If sharing is low → Don't invest in AT Protocol

**The 80/20 rule applies here:** Simple sharing gets you 80% of the value with 20% of the effort. AT Protocol gets you the last 20% but requires 80% more effort.

## What This Means for Users (If You Do Build AT Protocol)

### What Users Would Actually Be Able To Do

#### 1. **Share Their Meal Plans and Recipes on Bluesky**

**The Experience:**
- Sarah creates a meal plan for the week in Dinner Done Better
- She clicks a "Share to Bluesky" button
- A beautifully formatted post appears on her Bluesky feed showing:
  - The meal plan overview (e.g., "This week's dinners: Monday pasta, Tuesday tacos...")
  - Links back to view the full meal plan in Dinner Done Better
  - Photos of the meals (if she has them)
  - Her notes about the meal plan

**Why This Matters:**
- Users can share their cooking journey with friends and followers
- Meal plans become shareable content, not just private planning tools
- Friends can see what you're cooking and get inspired

#### 2. **Discover Recipes from the Community**

**The Experience:**
- John is looking for new dinner ideas
- He opens Dinner Done Better and sees a "Discover" tab
- He finds recipes that other users have shared on Bluesky
- He can see ratings, photos, and comments from the community
- With one click, he can add a discovered recipe to his own collection

**Why This Matters:**
- Users get access to a much larger recipe database
- Real people's tested recipes, not just curated content
- Social proof: "50 people made this and loved it"

#### 3. **Follow Other Meal Planners**

**The Experience:**
- Maria follows her friend Alex on Bluesky
- Alex shares his meal plans every week
- Maria sees Alex's meal plans in her Dinner Done Better feed
- She can "steal" his meal plans for her own week
- She can see what recipes Alex is using and try them herself

**Why This Matters:**
- Meal planning becomes collaborative and social
- Users can learn from each other
- Friends and families can coordinate meal planning

#### 4. **Build a Cooking Social Network**

**The Experience:**
- Users can follow other cooks on Bluesky who use Dinner Done Better
- They see a feed of meal plans and recipes from people they follow
- They can like, comment, and repost meal plans
- They can discover new cooks based on similar tastes

**Why This Matters:**
- Creates a community around meal planning
- Users find people with similar cooking styles
- Builds engagement and retention

#### 5. **Cross-Platform Recipe Sharing**

**The Experience:**
- A recipe shared on Bluesky can be opened directly in Dinner Done Better
- Users browsing Bluesky can discover Dinner Done Better through shared content
- Recipes become portable - they work across platforms

**Why This Matters:**
- Content discovery brings new users to the platform
- Users aren't locked into one platform
- Recipes become part of a larger ecosystem

### What Users Would NOT Be Able To Do (Initially)

- **Vote on meal plans via Bluesky** - Voting stays within Dinner Done Better
- **Edit recipes via Bluesky** - Recipe editing stays in the app
- **Full meal planning via Bluesky** - The core planning experience remains in Dinner Done Better
- **Private meal plans on Bluesky** - Only public shares go to Bluesky (privacy is preserved)

### The Value Proposition

**For Users:**
- **Social sharing**: Show off your meal planning skills
- **Discovery**: Find recipes from a community of real cooks
- **Inspiration**: See what others are cooking
- **Connection**: Follow friends and cooking influencers
- **Portability**: Your recipes aren't locked to one platform

**For Dinner Done Better:**
- **User acquisition**: Content discovery brings new users
- **Engagement**: Social features increase app usage
- **Retention**: Following and sharing creates stickiness
- **Network effects**: More users = more content = more value
- **Differentiation**: Unique social meal planning experience

### What This Actually Looks Like Day-to-Day

**Before AT Protocol Integration:**
- User creates meal plan → Stays private in the app
- User finds recipe → Has to manually search or import
- User wants to share → Takes screenshots and shares manually
- User wants inspiration → Limited to what's in the app

**After AT Protocol Integration:**
- User creates meal plan → One-click share to Bluesky, friends see it
- User finds recipe → Discovers it from Bluesky feed, adds with one click
- User wants to share → Automatic, beautiful posts with links
- User wants inspiration → Sees feed of meal plans from people they follow

**The Key Difference:**
Your meal planning app becomes **part of a social network** instead of an isolated tool. Every meal plan can become content. Every recipe can be discovered. Every user can become part of a cooking community.

### Is This Worth Building? (Updated Assessment)

**Comparison: Simple Sharing vs. AT Protocol Integration**

| Factor | Simple Sharing | AT Protocol Integration |
|-------|---------------|----------------------|
| **Development Time** | 1-2 days | 6-8 weeks |
| **Infrastructure** | None (just links) | OAuth, XRPC client, token storage |
| **User Friction** | Medium (copy/paste) | Low (one-click) |
| **Discovery Features** | ❌ No | ✅ Yes (can read Bluesky feed) |
| **Social Graph** | ❌ No | ✅ Yes (follow users in-app) |
| **Automatic Posting** | ❌ No | ✅ Yes |
| **Maintenance** | Low | Medium (API changes) |
| **User Acquisition** | Medium (manual sharing) | High (automatic + discovery) |

**When Simple Sharing Makes Sense:**
- ✅ You just want users to share content
- ✅ You're resource-constrained
- ✅ You want to validate demand first
- ✅ Your users are tech-savvy (don't mind copy/paste)
- ✅ You want to ship fast

**When AT Protocol Makes Sense:**
- ✅ You want discovery features (users find recipes from Bluesky)
- ✅ You want social graph integration (follow users in-app)
- ✅ You have 6-8 weeks to invest
- ✅ You're building a social platform, not just a tool
- ✅ You want automatic posting to reduce friction

**The Honest Recommendation:**

**Start Simple:**
1. Build share button that formats text + link (1-2 days)
2. Add public link views for meal plans/recipes (1 day)
3. Measure: Are users actually sharing?
4. If yes → Consider AT Protocol for automatic posting
5. If you want discovery → You'll need AT Protocol

**Don't build AT Protocol if:**
- You just want basic sharing (simple approach is fine)
- You're not sure users will share (validate first)
- You don't want discovery features (AT Protocol's main value)
- You're resource-constrained (simple approach gets 80% of value)

**Build AT Protocol if:**
- You want users to discover recipes from Bluesky (this is the killer feature)
- You want a social feed in your app
- You've validated that users want to share
- You have the resources and want to differentiate

**The Real Question:**
Do you want a **meal planning tool with sharing**, or a **social meal planning platform**? 

- Tool with sharing → Simple approach
- Social platform → AT Protocol

### What Simple Sharing Looks Like (Technically)

**What you'd build:**

```go
// Simple share formatter
func FormatMealPlanForSharing(mealPlan *MealPlan) string {
    var b strings.Builder
    b.WriteString("🍽️ My meal plan for this week:\n\n")
    for _, event := range mealPlan.Events {
        b.WriteString(fmt.Sprintf("• %s: %s\n", event.Date, event.Meal.Name))
    }
    b.WriteString(fmt.Sprintf("\nView full meal plan: %s/meal-plans/%s", 
        baseURL, mealPlan.ID))
    return b.String()
}

// API endpoint
func (s *Service) GetShareableMealPlan(ctx context.Context, mealPlanID string) (*ShareableMealPlan, error) {
    mealPlan, err := s.repo.GetMealPlan(ctx, mealPlanID)
    if err != nil {
        return nil, err
    }
    
    return &ShareableMealPlan{
        Text: FormatMealPlanForSharing(mealPlan),
        Link: fmt.Sprintf("%s/meal-plans/%s", baseURL, mealPlanID),
        QRCode: generateQRCode(mealPlanID),
    }, nil
}
```

**Frontend:**
- "Share" button that calls API
- Copies formatted text + link to clipboard
- Shows "Copied! Paste into Bluesky" message
- Maybe opens Bluesky web app with pre-filled text (if possible)

**Public link view:**
- `/meal-plans/{id}` route (public, no auth required)
- Shows meal plan details
- "Sign up to create your own" CTA

**That's it.** No OAuth, no XRPC, no token storage, no AT Protocol complexity.

**What you get:**
- Users can share meal plans
- Links drive traffic back to your app
- Public pages can be indexed by search engines
- Minimal development time
- No ongoing maintenance

**What you don't get:**
- Automatic posting (users still copy/paste)
- Discovery from Bluesky (can't read Bluesky feed)
- Social graph integration (can't follow users in-app)
- Rich embeds (links are just links)

**Is this enough?** For most use cases, yes. Only build AT Protocol if you specifically need discovery or social graph features.

## User Experience Examples

### Example 1: Sarah Shares Her Weekly Meal Plan

**Sarah's Flow:**
1. She creates a meal plan in Dinner Done Better for the week
2. She clicks "Share to Bluesky" (after connecting her Bluesky account once)
3. A preview shows what the post will look like
4. She adds a personal note: "Trying to meal prep better this week! Here's what I'm making."
5. The post goes live on Bluesky with:
   - A summary of her meal plan
   - Links to view details in Dinner Done Better
   - Photos of the meals
6. Her Bluesky followers see it in their feed
7. Some followers click through to Dinner Done Better to see the full meal plan

**What Happens Next:**
- Friends comment on Bluesky: "That pasta looks amazing!"
- Someone asks: "Can you share the recipe?"
- Sarah can share individual recipes from the meal plan
- New users discover Dinner Done Better through her post

### Example 2: John Discovers a New Recipe

**John's Flow:**
1. He's browsing his Bluesky feed
2. He sees a post from a cook he follows: "Made this amazing butter chicken last night! 🍗"
3. The post includes a link to the recipe in Dinner Done Better
4. He clicks through and sees the full recipe with ingredients, steps, and photos
5. He clicks "Add to My Recipes" in Dinner Done Better
6. The recipe is now in his collection
7. He can use it in his own meal plans

**What Happens Next:**
- John makes the recipe and loves it
- He shares his own meal plan featuring that recipe
- The original cook sees John's post and feels good about sharing
- More people discover the recipe through John's share

### Example 3: Maria Follows Her Friend's Meal Planning

**Maria's Flow:**
1. She connects her Bluesky account to Dinner Done Better
2. She follows her friend Alex on Bluesky (who also uses Dinner Done Better)
3. In her Dinner Done Better app, she sees a new "Following" feed
4. She sees Alex's meal plan for the week
5. She thinks "That looks good!" and clicks "Use This Meal Plan"
6. The meal plan is copied to her account (she can modify it)
7. All the recipes are automatically added to her collection

**What Happens Next:**
- Maria and Alex coordinate their meal planning
- They share tips and modifications
- Their families eat similar meals and can discuss them
- They discover new recipes together

### Example 4: A Cooking Influencer Builds an Audience

**The Influencer's Flow:**
1. A cooking influencer uses Dinner Done Better for meal planning
2. They share their meal plans weekly on Bluesky
3. Their followers see beautiful, well-planned meals
4. Followers click through to Dinner Done Better to get the recipes
5. The influencer gains followers who want their meal plans
6. They can monetize by sharing premium meal plans

**What Happens Next:**
- The influencer's content drives new users to Dinner Done Better
- Their followers become Dinner Done Better users
- Network effects: more users = more content = more value
- Dinner Done Better becomes the platform for social meal planning

## Why Integrate AT Protocol?

### The Business Case

**User Acquisition:**
- Content shared on Bluesky drives discovery
- Each shared meal plan is a potential new user
- Network effects: more users = more content = more users

**User Engagement:**
- Social features increase daily active users
- Following and sharing creates habit loops
- Users return to see what their network is cooking

**User Retention:**
- Social connections create stickiness
- Users invested in their social graph won't leave easily
- Community building increases lifetime value

**Differentiation:**
- No other meal planning app has deep Bluesky integration
- Creates a unique value proposition
- Positions Dinner Done Better as the social meal planning platform

### The Technical Benefits

- **No Infrastructure Overhead**: Use Bluesky's existing servers (for MVP)
- **Proven Protocol**: AT Protocol is battle-tested by Bluesky
- **Future-Proof**: Decentralized means not locked to one provider
- **Interoperability**: Works with any AT Protocol-compatible service

## Integration Approaches

### Option 1: Client-Only Integration (Recommended for MVP)

**What**: Your backend acts as an AT Protocol client, connecting to existing PDS (Personal Data Server) infrastructure (like Bluesky's).

**Pros**:
- Fastest to implement
- No infrastructure overhead
- Leverages existing, stable services
- Good for MVP and testing

**Cons**:
- Dependency on external services
- Limited customization
- Less control over user data storage

**Use Case**: Allow users to link their Bluesky accounts and share meal plans/recipes to their Bluesky feed.

### Option 2: Hybrid Approach

**What**: Your backend acts as both a client (for reading/sharing) and optionally hosts a lightweight PDS for power users.

**Pros**:
- Flexibility for different user needs
- Can start simple and expand
- Some users can have full control

**Cons**:
- More complex architecture
- Need to manage PDS infrastructure

**Use Case**: Most users use Bluesky's PDS, but power users can optionally host their own.

### Option 3: Full PDS Implementation

**What**: Run your own Personal Data Server for all users.

**Pros**:
- Complete control
- Custom schemas and features
- No external dependencies

**Cons**:
- Significant infrastructure complexity
- Storage and bandwidth requirements
- Maintenance overhead
- Not recommended for initial implementation

**Recommendation**: Start with **Option 1** (Client-Only) for MVP, then consider Option 2 if needed.

## Architecture Design

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Dinner Done Better                       │
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐       │
│  │   Identity   │  │ Meal Planning│  │   Recipes    │       │
│  │   Domain     │  │    Domain    │  │    Domain    │       │
│  └──────────────┘  └──────────────┘  └──────────────┘  x     │
│         │                  │                  │             │
│         └──────────────────┼──────────────────┘             │
│                            │                                │
│                   ┌────────▼────────┐                       │
│                   │  AT Protocol    │                       │
│                   │  Integration    │                       │
│                   │    Service      │                       │
│                   └────────┬────────┘                       │
│                            │                                │
└────────────────────────────┼────────────────────────────────┘
                             │
                             │ XRPC / HTTP
                             │
                ┌────────────▼────────────┐
                │   AT Protocol Network   │
                │  (Bluesky PDS, etc.)    │
                └─────────────────────────┘
```

### Domain Structure

Create a new domain package: `internal/domain/atprotocol/`

```
internal/domain/atprotocol/
├── atprotocol.go              # Core domain models
├── identity_link.go           # Links between DDB users and AT Protocol identities
├── share.go                   # Shareable content (meal plans, recipes)
├── repository.go              # Repository interface
└── wire.go                    # Wire dependency injection
```

## Implementation Steps

### Phase 1: Foundation (Week 1-2)

1. **Add Go AT Protocol Client Library**
   ```bash
   go get github.com/bluesky-social/indigo
   ```

2. **Create Domain Models**
   - `ATProtocolIdentityLink`: Links a DDB user to an AT Protocol identity (DID/handle)
   - `ShareableContent`: Represents content that can be shared (meal plans, recipes)
   - `Share`: Represents a share action with metadata

3. **Database Schema**
   - Create `at_protocol_identity_links` table
   - Create `at_protocol_shares` table
   - Add indexes for lookups

4. **Repository Layer**
   - Implement repository interfaces for AT Protocol domain
   - Use existing patterns from other domain packages

### Phase 2: Authentication & Identity Linking (Week 3-4)

1. **OAuth Flow**
   - Implement AT Protocol OAuth 2.1 flow
   - Create endpoints for:
     - Initiate OAuth flow
     - OAuth callback handler
     - Link/unlink AT Protocol identity

2. **Identity Resolution**
   - Resolve handles to DIDs
   - Store DID and handle in identity link
   - Handle handle changes

3. **User Interface**
   - Add "Connect Bluesky" button in user settings
   - Show linked identity status
   - Allow unlinking

### Phase 3: Content Sharing (Week 5-6)

1. **Share Service**
   - Create service to format meal plans/recipes for AT Protocol
   - Generate rich text posts with recipe details
   - Handle media (recipe images)

2. **XRPC Client**
   - Implement client for creating posts
   - Handle authentication tokens
   - Error handling and retries

3. **Share Endpoints**
   - `POST /api/v1/meal-plans/{id}/share` - Share meal plan
   - `POST /api/v1/recipes/{id}/share` - Share recipe
   - `GET /api/v1/shares` - List user's shares

### Phase 4: Reading & Discovery (Week 7-8)

1. **Feed Integration**
   - Subscribe to AT Protocol firehose (optional)
   - Parse meal planning-related posts
   - Extract recipe/meal plan data from posts

2. **Discovery Features**
   - Search for recipes shared on AT Protocol
   - Follow other meal planners
   - Display social feed of followed users

### Phase 5: Advanced Features (Future)

1. **Custom Lexicons**
   - Define custom schemas for meal plans and recipes
   - Enable structured data sharing
   - Better discovery and parsing

2. **Social Graph**
   - Follow/unfollow users
   - Like and repost meal plans
   - Comments and discussions

## Go Libraries and Tools

### Primary Library

**github.com/bluesky-social/indigo**

This is the official Go implementation of the AT Protocol. It provides:
- XRPC client/server
- Lexicon definitions
- DID resolution
- Repository operations
- Authentication

```go
import (
    "github.com/bluesky-social/indigo/atproto"
    "github.com/bluesky-social/indigo/xrpc"
)
```

### Additional Libraries

1. **DID Resolution**: `github.com/multiformats/go-multibase` (for DID handling)
2. **OAuth**: Use existing `github.com/go-oauth2/oauth2/v4` or AT Protocol's OAuth implementation
3. **JSON-LD**: For structured data (if using custom lexicons)

## Domain Models

### ATProtocolIdentityLink

```go
package atprotocol

import "time"

type ATProtocolIdentityLink struct {
    ID                string
    BelongsToUser     string  // DDB User ID
    DID               string  // Decentralized Identifier
    Handle            string  // AT Protocol handle (e.g., @user.bsky.social)
    AccessToken       string  // OAuth access token (encrypted)
    RefreshToken      string  // OAuth refresh token (encrypted)
    TokenExpiresAt    *time.Time
    CreatedAt         time.Time
    LastUpdatedAt     *time.Time
    ArchivedAt        *time.Time
}
```

### ShareableContent

```go
type ShareableContentType string

const (
    ShareableContentTypeMealPlan ShareableContentType = "meal_plan"
    ShareableContentTypeRecipe   ShareableContentType = "recipe"
)

type ShareableContent struct {
    ContentType ShareableContentType
    ContentID   string  // ID of the meal plan or recipe
    Title       string
    Description string
    ImageURL    *string
    ShareURL    string  // URL to view in DDB
}
```

### Share

```go
type Share struct {
    ID                string
    BelongsToUser     string
    ATProtocolPostURI string  // AT URI of the created post (at://...)
    ContentType       ShareableContentType
    ContentID         string
    CreatedAt         time.Time
    ArchivedAt        *time.Time
}
```

## API Design

### gRPC Service

Add to `proto/atprotocol/atprotocol_service.proto`:

```protobuf
syntax = "proto3";

package atprotocol;

import "atprotocol/atprotocol_service_types.proto";

service ATProtocolService {
  // Link a user's AT Protocol identity
  rpc LinkIdentity(LinkIdentityRequest) returns (LinkIdentityResponse);
  
  // Unlink a user's AT Protocol identity
  rpc UnlinkIdentity(UnlinkIdentityRequest) returns (UnlinkIdentityResponse);
  
  // Get linked identity status
  rpc GetLinkedIdentity(GetLinkedIdentityRequest) returns (GetLinkedIdentityResponse);
  
  // Share a meal plan to AT Protocol
  rpc ShareMealPlan(ShareMealPlanRequest) returns (ShareMealPlanResponse);
  
  // Share a recipe to AT Protocol
  rpc ShareRecipe(ShareRecipeRequest) returns (ShareRecipeResponse);
  
  // List user's shares
  rpc ListShares(ListSharesRequest) returns (ListSharesResponse);
  
  // Delete a share (archive the post)
  rpc DeleteShare(DeleteShareRequest) returns (DeleteShareResponse);
}
```

### HTTP Endpoints

For OAuth flow and web UI:

```go
// OAuth initiation
GET  /api/v1/atprotocol/oauth/authorize
GET  /api/v1/atprotocol/oauth/callback

// Identity management
GET  /api/v1/atprotocol/identity
POST /api/v1/atprotocol/identity/link
DELETE /api/v1/atprotocol/identity/unlink

// Sharing
POST /api/v1/atprotocol/meal-plans/{id}/share
POST /api/v1/atprotocol/recipes/{id}/share
GET  /api/v1/atprotocol/shares
DELETE /api/v1/atprotocol/shares/{id}
```

## Security Considerations

1. **Token Storage**
   - Encrypt OAuth tokens at rest
   - Use your existing encryption utilities from `internal/platform`
   - Never log tokens

2. **Authentication**
   - Verify user owns the AT Protocol identity before linking
   - Use secure OAuth flow with PKCE
   - Validate tokens on each API call

3. **Rate Limiting**
   - Respect AT Protocol rate limits
   - Implement backoff for retries
   - Queue shares if rate limited

4. **Content Validation**
   - Validate content before sharing
   - Sanitize user-generated content
   - Respect user privacy settings

5. **Error Handling**
   - Don't expose internal errors to users
   - Log errors for debugging
   - Handle network failures gracefully

## Database Schema

### at_protocol_identity_links

```sql
CREATE TABLE at_protocol_identity_links (
    id TEXT PRIMARY KEY,
    belongs_to_user TEXT NOT NULL REFERENCES users(id),
    did TEXT NOT NULL UNIQUE,
    handle TEXT NOT NULL,
    access_token_encrypted TEXT NOT NULL,
    refresh_token_encrypted TEXT,
    token_expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_updated_at TIMESTAMPTZ,
    archived_at TIMESTAMPTZ,
    
    CONSTRAINT fk_user FOREIGN KEY (belongs_to_user) REFERENCES users(id)
);

CREATE INDEX idx_at_protocol_identity_links_user ON at_protocol_identity_links(belongs_to_user) WHERE archived_at IS NULL;
CREATE INDEX idx_at_protocol_identity_links_did ON at_protocol_identity_links(did) WHERE archived_at IS NULL;
```

### at_protocol_shares

```sql
CREATE TABLE at_protocol_shares (
    id TEXT PRIMARY KEY,
    belongs_to_user TEXT NOT NULL REFERENCES users(id),
    at_protocol_post_uri TEXT NOT NULL UNIQUE,
    content_type TEXT NOT NULL, -- 'meal_plan' or 'recipe'
    content_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMPTZ,
    
    CONSTRAINT fk_user FOREIGN KEY (belongs_to_user) REFERENCES users(id)
);

CREATE INDEX idx_at_protocol_shares_user ON at_protocol_shares(belongs_to_user) WHERE archived_at IS NULL;
CREATE INDEX idx_at_protocol_shares_content ON at_protocol_shares(content_type, content_id) WHERE archived_at IS NULL;
```

## Example Implementation

### Service Layer

```go
package atprotocol

import (
    "context"
    "github.com/bluesky-social/indigo/xrpc"
)

type Service interface {
    LinkIdentity(ctx context.Context, userID string, authCode string) error
    ShareMealPlan(ctx context.Context, userID string, mealPlanID string) (*Share, error)
    ShareRecipe(ctx context.Context, userID string, recipeID string) (*Share, error)
}

type service struct {
    repo          Repository
    identityRepo  identity.Repository
    mealPlanRepo  mealplanning.Repository
    recipeRepo    mealplanning.RecipeRepository
    xrpcClient    *xrpc.Client
    logger        logging.Logger
}

func (s *service) ShareMealPlan(ctx context.Context, userID string, mealPlanID string) (*Share, error) {
    // 1. Get user's linked AT Protocol identity
    link, err := s.repo.GetIdentityLinkByUserID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // 2. Get meal plan
    mealPlan, err := s.mealPlanRepo.GetMealPlan(ctx, mealPlanID)
    if err != nil {
        return nil, err
    }
    
    // 3. Format content for AT Protocol
    text := formatMealPlanForPost(mealPlan)
    
    // 4. Create post via XRPC
    postURI, err := s.createPost(ctx, link, text)
    if err != nil {
        return nil, err
    }
    
    // 5. Store share record
    share := &Share{
        BelongsToUser:     userID,
        ATProtocolPostURI: postURI,
        ContentType:       ShareableContentTypeMealPlan,
        ContentID:         mealPlanID,
    }
    
    return s.repo.CreateShare(ctx, share)
}

func (s *service) createPost(ctx context.Context, link *ATProtocolIdentityLink, text string) (string, error) {
    client := &xrpc.Client{
        Client: s.xrpcClient,
        Auth: &xrpc.AuthInfo{
            AccessJwt: link.AccessToken,
        },
    }
    
    // Use indigo library to create post
    // Implementation depends on specific AT Protocol client library API
    // ...
    
    return postURI, nil
}
```

## Testing Strategy

1. **Unit Tests**
   - Test domain models and business logic
   - Mock XRPC client
   - Test error handling

2. **Integration Tests**
   - Test OAuth flow with test PDS
   - Test sharing functionality
   - Test identity linking/unlinking

3. **E2E Tests**
   - Full flow: link identity → share content → verify post
   - Test with Bluesky test environment

## Next Steps

1. Review and approve this integration plan
2. Set up development environment with AT Protocol test PDS
3. Create domain models and database schema
4. Implement OAuth flow
5. Build sharing functionality
6. Add UI components for linking and sharing
7. Test with real Bluesky accounts (test environment)
8. Deploy to staging for user testing

## Resources

- [AT Protocol Specification](https://atproto.com/specs/atp)
- [Bluesky AT Protocol GitHub](https://github.com/bluesky-social/atproto)
- [Indigo (Go SDK) GitHub](https://github.com/bluesky-social/indigo)
- [AT Protocol Guide](https://atproto.com/guides)
- [OAuth Guide](https://atproto.com/guides/oauth)
