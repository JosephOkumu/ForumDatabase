// Global state
let currentUser = null;

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    loadPosts(); // Load posts when page loads
});

// UI Helper Functions
function showLoginForm() {
    document.getElementById('loginForm').style.display = 'block';
    document.getElementById('registerForm').style.display = 'none';
    document.getElementById('createPostForm').style.display = 'none';
}

function showRegisterForm() {
    document.getElementById('loginForm').style.display = 'none';
    document.getElementById('registerForm').style.display = 'block';
    document.getElementById('createPostForm').style.display = 'none';
}

function showCreatePostForm() {
    console.log('Current user:', currentUser); // Debug log
    if (!currentUser) {
        alert('Please login to create a post');
        return;
    }
    document.getElementById('loginForm').style.display = 'none';
    document.getElementById('registerForm').style.display = 'none';
    document.getElementById('createPostForm').style.display = 'block';
}

// Update the handleCreatePost function
async function handleCreatePost(event) {
    event.preventDefault();
    if (!currentUser) {
        alert('Please login to create a post');
        return;
    }

    const title = document.getElementById('postTitle').value;
    const content = document.getElementById('postContent').value;
    const categoriesSelect = document.getElementById('postCategories');
    const categories = Array.from(categoriesSelect.selectedOptions).map(option => option.value);

    try {
        const response = await fetch('/create-post', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                title,
                content,
                categories,
                userId: currentUser.id // Include user ID if needed
            }),
            credentials: 'include' // Important for sending cookies
        });

        if (response.ok) {
            document.getElementById('createPostForm').style.display = 'none';
            document.getElementById('postTitle').value = '';
            document.getElementById('postContent').value = '';
            categoriesSelect.selectedIndex = -1;
            alert('Post created successfully!');
            loadPosts();
        } else {
            const errorData = await response.json();
            alert(errorData.error || 'Failed to create post.');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('An error occurred while creating the post.');
    }
}

// Filter functions
function filterPosts() {
    const category = document.getElementById('categoryFilter').value;
    if (category) {
        fetch(`/category?category=${encodeURIComponent(category)}`)
            .then(response => response.json())
            .then(posts => displayPosts(posts))
            .catch(error => {
                console.error('Error:', error);
                alert('Failed to filter posts');
            });
    } else {
        loadPosts();
    }
}

function showCreatedPosts() {
    if (!currentUser) {
        alert('Please login to view your posts');
        return;
    }
    // Add logic to fetch and display user's created posts
}

function showLikedPosts() {
    if (!currentUser) {
        alert('Please login to view liked posts');
        return;
    }
    // Add logic to fetch and display user's liked posts
}

// API Functions
async function handleLogin(event) {
    event.preventDefault();
    const email = document.getElementById('loginEmail').value;
    const password = document.getElementById('loginPassword').value;

    try {
        const response = await fetch('/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, password }),
            credentials: 'include' // This is important for handling cookies
        });

        if (response.ok) {
            const data = await response.json();
            currentUser = data; // Store the entire user data
            document.getElementById('loginForm').style.display = 'none';
            alert('Login successful!');
            loadPosts();
            updateUIForLoggedInUser();
        } else {
            const errorData = await response.json();
            alert(errorData.error || 'Login failed. Please check your credentials.');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('An error occurred during login.');
    }
}

async function checkLoginStatus() {
    try {
        const response = await fetch('/check-session', {
            credentials: 'include'
        });
        if (response.ok) {
            const data = await response.json();
            currentUser = data;
            updateUIForLoggedInUser();
        }
    } catch (error) {
        console.error('Error checking login status:', error);
    }
}

function updateUIForLoggedInUser() {
    const navLinks = document.querySelector('.nav-links');
    if (currentUser) {
        navLinks.innerHTML = `
            <span>Welcome, ${currentUser.username}</span>
            <button class="btn" onclick="showCreatePostForm()">Create Post</button>
            <button class="btn" onclick="handleLogout()">Logout</button>
        `;
    }
}

async function handleRegister(event) {
    event.preventDefault();
    const email = document.getElementById('registerEmail').value;
    const username = document.getElementById('registerUsername').value;
    const password = document.getElementById('registerPassword').value;

    try {
        const response = await fetch('/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, username, password }),
        });

        if (response.ok) {
            alert('Registration successful! Please login.');
            showLoginForm();
        } else {
            const errorData = await response.json();
            alert(errorData.error || 'Registration failed. Please try again.');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('An error occurred during registration.');
    }
}

async function handleCreatePost(event) {
    event.preventDefault();
    if (!currentUser) {
        alert('Please login to create a post');
        return;
    }

    const title = document.getElementById('postTitle').value;
    const content = document.getElementById('postContent').value;
    const categoriesSelect = document.getElementById('postCategories');
    const categories = Array.from(categoriesSelect.selectedOptions).map(option => option.value);

    try {
        const response = await fetch('/create-post', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ title, content, categories }),
        });

        if (response.ok) {
            document.getElementById('createPostForm').style.display = 'none';
            document.getElementById('postTitle').value = '';
            document.getElementById('postContent').value = '';
            categoriesSelect.selectedIndex = -1;
            loadPosts();
            alert('Post created successfully!');
        } else {
            const errorData = await response.json();
            alert(errorData.error || 'Failed to create post.');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('An error occurred while creating the post.');
    }
}

async function loadPosts() {
    try {
        const response = await fetch('/posts');
        if (!response.ok) {
            throw new Error('Failed to fetch posts');
        }
        const posts = await response.json();
        displayPosts(posts);
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to load posts.');
    }
}

function displayPosts(posts) {
    const container = document.getElementById('postsContainer');
    if (!container) return;

    container.innerHTML = '';
    if (!Array.isArray(posts) || posts.length === 0) {
        container.innerHTML = '<p>No posts available.</p>';
        return;
    }

    posts.forEach(post => {
        const postElement = document.createElement('div');
        postElement.className = 'post';
        postElement.innerHTML = `
            <div class="post-header">
                <h3>${escapeHtml(post.title)}</h3>
                <span>Posted by ${escapeHtml(post.username)}</span>
            </div>
            <div class="post-categories">
                ${post.categories ? post.categories.map(category => 
                    `<span class="category-tag">${escapeHtml(category)}</span>`
                ).join('') : ''}
            </div>
            <p>${escapeHtml(post.content)}</p>
            <div class="reaction-buttons">
                <button class="reaction-btn" onclick="handleReaction(${post.id}, 'like')">
                    👍 <span>${post.likes || 0}</span>
                </button>
                <button class="reaction-btn" onclick="handleReaction(${post.id}, 'dislike')">
                    👎 <span>${post.dislikes || 0}</span>
                </button>
            </div>
            <div class="comments" id="comments-${post.id}">
                ${renderComments(post.comments)}
                ${currentUser ? `
                    <form onsubmit="handleAddComment(event, ${post.id})" class="comment-form">
                        <div class="form-group">
                            <textarea required placeholder="Add a comment..." class="comment-input"></textarea>
                        </div>
                        <button type="submit" class="btn btn-secondary">Comment</button>
                    </form>
                ` : ''}
            </div>
        `;
        container.appendChild(postElement);
    });
}

function escapeHtml(unsafe) {
    if (!unsafe) return '';
    return unsafe
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}

function renderComments(comments) {
    if (!comments || !Array.isArray(comments) || comments.length === 0) return '';
    
    return comments.map(comment => `
        <div class="comment">
            <p>${escapeHtml(comment.content)}</p>
            <small>By ${escapeHtml(comment.username)}</small>
            <div class="reaction-buttons">
                <button class="reaction-btn" onclick="handleCommentReaction(${comment.id}, 'like')">
                    👍 <span>${comment.likes || 0}</span>
                </button>
                <button class="reaction-btn" onclick="handleCommentReaction(${comment.id}, 'dislike')">
                    👎 <span>${comment.dislikes || 0}</span>
                </button>
            </div>
        </div>
    `).join('');
}

async function handleReaction(postId, type) {
    if (!currentUser) {
        alert('Please login to react to posts');
        return;
    }

    try {
        const response = await fetch('/add-reaction', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ postId, type }),
        });

        if (!response.ok) {
            throw new Error('Failed to add reaction');
        }
        loadPosts();
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to add reaction.');
    }
}

async function handleCommentReaction(commentId, type) {
    if (!currentUser) {
        alert('Please login to react to comments');
        return;
    }

    try {
        const response = await fetch('/commentreaction', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ commentId, type }),
        });

        if (!response.ok) {
            throw new Error('Failed to add comment reaction');
        }
        loadPosts();
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to add comment reaction.');
    }
}

async function handleAddComment(event, postId) {
    event.preventDefault();
    if (!currentUser) {
        alert('Please login to comment');
        return;
    }

    const textarea = event.target.querySelector('textarea');
    const content = textarea.value.trim();

    if (!content) {
        alert('Please enter a comment');
        return;
    }

    try {
        const response = await fetch('/comment', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ postId, content }),
        });

        if (!response.ok) {
            throw new Error('Failed to add comment');
        }

        textarea.value = '';
        loadPosts();
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to add comment.');
    }
}

function handleLogout() {
    currentUser = null;
    updateUIForLoggedInUser();
    loadPosts();
}

// Update the DOMContentLoaded event listener
document.addEventListener('DOMContentLoaded', function() {
    checkLoginStatus(); // Check if user is already logged in
    loadPosts();
});