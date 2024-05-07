# Project Requirements

## Stage 1: Authentication and Profile Creation

- Implement authentication system. [x]
- Allow users to create profiles. [x]
- Need to make a function for universal hashing in order to be used for password and paths 

## Stage 2: Backend System for Folder Creation

- Allow particular users to set up folders via a method and store their paths along with a name for that path.
- Hash paths and then make a QR out of it
- Develop a backend system to create folders.
- Store folder information in a database or filesystem.
- Ensure folders are accessible to users.

## Stage 3: QR Code Sharing for Read-Only Access

- Generate QR codes for folder access.
- Verify user profile or permissions before granting read access.
- Prioritize functionality over sign-ups.

## Stage 4: Simple UI Implementation

- Create a simple user interface.
- Include components for user profiles, folder creation, and subscribed folders.

## Stage 5: Adding Files to QR-Enabled Directories

- Enable users to add files to directories accessed via QR codes.
- Allow access without signing in, based on QR code authentication.
- Consider implementing access based on user credentials as an option.

## Payment Integration

- Implement a payment system for storage usage.
- Optionally, set up local storage to minimize costs.
