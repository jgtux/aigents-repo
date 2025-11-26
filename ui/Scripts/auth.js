// auth.js - Authentication module for AIgents
// Place this file in: ui/Scripts/auth.js

const API_BASE_URL = 'http://localhost:8080/api/v1';

/**
 * Sign up a new user
 * @param {string} email - User's email
 * @param {string} password - User's password (min 8, max 25 chars)
 * @returns {Promise<{success: boolean, data?: any, error?: string}>}
 */
async function signUp(email, password) {
  try {
    // Validate input
    if (!email || !password) {
      return { success: false, error: 'Email and password are required' };
    }

    if (password.length < 8 || password.length > 25) {
      return { success: false, error: 'Password must be between 8 and 25 characters' };
    }

    const response = await fetch(`${API_BASE_URL}/auth/create`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include', // Important for HTTP-only cookies
      body: JSON.stringify({
        email: email,
        password: password
      })
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error || data.message || 'Signup failed');
    }

    return { success: true, data };

  } catch (error) {
    console.error('Signup error:', error);
    return { success: false, error: error.message };
  }
}

/**
 * Log in an existing user
 * @param {string} email - User's email
 * @param {string} password - User's password
 * @returns {Promise<{success: boolean, data?: any, error?: string}>}
 */
async function login(email, password) {
  try {
    // Validate input
    if (!email || !password) {
      return { success: false, error: 'Email and password are required' };
    }

    const response = await fetch(`${API_BASE_URL}/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include', // Important for HTTP-only cookies
      body: JSON.stringify({
        email: email,
        password: password
      })
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error || data.message || 'Login failed');
    }

    // Cookies are automatically set by the browser
    return { success: true, data };

  } catch (error) {
    console.error('Login error:', error);
    return { success: false, error: error.message };
  }
}

/**
 * Refresh the access token using the refresh token
 * @returns {Promise<{success: boolean, data?: any, error?: string}>}
 */
async function refreshToken() {
  try {
    const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
      method: 'GET',
      credentials: 'include' // Important for HTTP-only cookies
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error || data.message || 'Token refresh failed');
    }

    return { success: true, data };

  } catch (error) {
    console.error('Token refresh error:', error);
    return { success: false, error: error.message };
  }
}

/**
 * Log out the current user
 * @returns {Promise<{success: boolean}>}
 */
async function logout() {
  try {
    // Note: Your backend doesn't have a logout endpoint yet
    // You may need to add one, or just clear local state
    
    // If you add a logout endpoint later:
    // await fetch(`${API_BASE_URL}/auth/logout`, {
    //   method: 'POST',
    //   credentials: 'include'
    // });

    // Redirect to login page
    window.location.href = '/Index.html/Login.html';
    
    return { success: true };
  } catch (error) {
    console.error('Logout error:', error);
    return { success: false };
  }
}

/**
 * Make an authenticated API request
 * @param {string} endpoint - API endpoint (without base URL)
 * @param {object} options - Fetch options
 * @returns {Promise<Response>}
 */
async function authenticatedRequest(endpoint, options = {}) {
  const url = endpoint.startsWith('http') ? endpoint : `${API_BASE_URL}${endpoint}`;
  
  const response = await fetch(url, {
    ...options,
    credentials: 'include', // Always include credentials for auth
    headers: {
      'Content-Type': 'application/json',
      ...options.headers
    }
  });

  // If unauthorized, try to refresh token once
  if (response.status === 401) {
    const refreshResult = await refreshToken();
    
    if (refreshResult.success) {
      // Retry the original request
      return fetch(url, {
        ...options,
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
          ...options.headers
        }
      });
    } else {
      // Redirect to login if refresh fails
      window.location.href = '/Index.html/Login.html';
    }
  }

  return response;
}

// Export functions if using modules, otherwise they're globally available
if (typeof module !== 'undefined' && module.exports) {
  module.exports = {
    signUp,
    login,
    refreshToken,
    logout,
    authenticatedRequest
  };
}