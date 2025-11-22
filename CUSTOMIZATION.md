# Customization Guide

This guide explains how to customize the application for your corporation, squad, or group.

## Corporation Name Branding

The application supports dynamic branding through the `CORPORATION_NAME` environment variable. This allows you to easily rebrand the entire application without modifying code.

### Configuration

#### Docker Deployment

**Option 1: Set via docker-compose command** (Recommended)
```bash
CORPORATION_NAME="Your Corp Name" docker-compose up -d
```

**Option 2: Create a `.env` file in the repository root**
```env
CORPORATION_NAME=Your Corp Name
```

Then start the containers:
```bash
docker-compose down
docker-compose up -d
```

The corporation name will be automatically applied to:
- Navigation header
- Page titles
- Admin panel branding
- Login screens
- User interface elements

#### Local Development

For local development, set the environment variable in `frontend/.env.local`:

```env
NEXT_PUBLIC_CORPORATION_NAME=Your Corp Name
```

Then restart your development server:
```bash
cd frontend
npm run dev
```

### Default Value

If `CORPORATION_NAME` is not set, the application defaults to "Tundragon Corp".

### What Gets Updated

When you change the `CORPORATION_NAME` variable, the following elements automatically update:

- **Navigation Menu**: The main header displays your corporation name
- **Page Title**: Browser tab title
- **Admin Panel**:
  - Welcome message
  - Sidebar branding
  - Login page description
- **Login Forms**: All authentication screens reference your corporation name

### Example Customizations

#### For a Corporation
```env
CORPORATION_NAME=Goonswarm Federation
```

#### For a Squad
```env
CORPORATION_NAME=Elite PVP Squad
```

#### For an Alliance
```env
CORPORATION_NAME=Test Alliance
```

### Technical Details

The application uses a simple configuration system:

1. **Docker Compose**: Passes `CORPORATION_NAME` environment variable to frontend container as `NEXT_PUBLIC_CORPORATION_NAME`
2. **Frontend**: React context (`AppConfigContext`) reads from `process.env.NEXT_PUBLIC_CORPORATION_NAME` and provides it to all components
3. **No Backend Dependencies**: Branding is frontend-only, no API calls needed

### Troubleshooting

If the corporation name doesn't update:

1. **Docker Deployment:**
   - Verify `CORPORATION_NAME` is set (either via command line or `.env` file)
   - Rebuild and restart: `docker-compose up -d --build`
   - Check the environment was passed: `docker exec eve_frontend env | grep CORPORATION`

2. **Local Development:**
   - Verify `NEXT_PUBLIC_CORPORATION_NAME` is set in `frontend/.env.local`
   - Restart the Next.js dev server completely (stop and start, not just refresh)
   - Check browser console for any errors

3. **General:**
   - Clear browser cache if the old name persists in page title
   - Remember: Environment variables starting with `NEXT_PUBLIC_` are baked into the build, so rebuilding may be necessary

### Future Customization Options

Potential future enhancements:
- Custom logo upload
- Color theme customization
- Configurable alliance/corporation ticker
- Custom footer text
