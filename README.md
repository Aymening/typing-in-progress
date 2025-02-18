# Web Forum Project

A web-based discussion platform that allows users to communicate, share posts, and interact with content through a simple and intuitive interface.

## Project Overview

This forum is designed to facilitate user communication and content sharing with essential features for a modern web forum experience.

## Key Features

### User Management
- User registration with email verification
- Secure login system with session management
- Single active session per user using cookies

### Post and Comments
- Create and view discussion posts
- Comment on existing posts
- Multiple category assignments for posts
- Public access for viewing content
- Posting and commenting limited to registered users

### Content Categories
- Organize posts by categories
- Filter posts based on categories
- Custom category creation and assignment

### Interaction System
- Like and dislike functionality for posts and comments
- Visible interaction counts for all users
- Interaction features limited to registered users

### Content Filtering
- Filter posts by categories (all users)
- View posts created by the logged-in user
- Access liked posts by the current user

## Technical Overview

### Database

![database](formatDatabase.png)

The project uses SQLite for data storage and management, providing:
- Efficient data storage
- SQL query support
- Secure user data handling

### Authentication
- Email-based registration
- Password encryption
- Session management with cookies
- Unique email verification

## Target Users

- General users (unregistered): Can view posts, comments, and interactions
- Registered users: Can create posts, comment, like/dislike, and filter personal content
- All users: Can filter content by categories


The project aims to create an engaging forum environment where users can easily share thoughts, interact with content, and participate in community discussions.

```

├── go.mod
├── go.sum
├── main.go
├── README.md
|
├── server
│   ├── database
│   │   └── database.db
|   |
│   ├── forum
│   │   ├── crud.go
│   │   ├── database.go
│   │   └── utils.go
|   |
│   └── handler
│       ├── api_login_signin.go
│       ├── Categorie.go
│       ├── comment.go
│       ├── errorhandler.go
│       ├── filter_post.go
│       ├── handlers.go
│       ├── isLoggedIn.go
│       ├── posts_section.go
│       ├── reactions.go
│       ├── UserPosts.go
│       └── UserProfile.go
|
└── website
    ├── css
    │   ├── err.css
    │   ├── nav.css
    │   ├── post.css
    │   ├── profile.css
    │   ├── style.css
    │   └── styleLogin.css
    |
    ├── img
    │   ├── cover.jpg
    │   ├── default.avif
    │   ├── favicon.ico
    │   └── logo.png
    |
    ├── js
    │   ├── auth.js
    │   ├── categ.js
    │   ├── comment.js
    │   ├── likePostLogic.js
    │   ├── likes_section.js
    │   ├── post.js
    │   ├── posts_section.js
    │   ├── profile_post.js
    │   ├── profil_section.js
    │   ├── script.js
    │   └── userSignLogin.js
    |
    └── pages
        ├── error1.html
        ├── index.html
        ├── login.html
        └── profil.html

```