import React, { useState, useEffect } from 'react';

function App() {
  const [tweets, setTweets] = useState([]);
  const [text, setText] = useState('');
  const [username, setUsername] = useState('');
  const [editingId, setEditingId] = useState(null);
  const [editingText, setEditingText] = useState('');

  // Load tweets
  const loadTweets = () => {
    fetch("http://localhost:8080/tweets")
      .then(res => res.json())
      .then(data => setTweets(data ?? []));
  };

  useEffect(() => {
    loadTweets();
  }, []);

  const handlePost = async () => {
    if (!text.trim() || !username.trim()) return;

    await fetch("http://localhost:8080/tweets", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, text })
    });

    setText('');
    loadTweets();
  };

  const handleDelete = async (id) => {
    await fetch(`http://localhost:8080/tweets/${id}`, {
      method: "DELETE"
    });
    loadTweets();
  };

  const startEditing = (id, currentText) => {
    setEditingId(id);
    setEditingText(currentText);
  };

  const handleEditSave = async () => {
    await fetch(`http://localhost:8080/tweets/${editingId}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ text: editingText })
    });
    setEditingId(null);
    setEditingText('');
    loadTweets();
  };

  return (
    <div style={{ padding: "2rem", maxWidth: "600px", margin: "auto" }}>
      <h1>🕊️ Twitter Clone</h1>

      <input
        type="text"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
        placeholder="Username"
        style={{ width: "100%", marginBottom: "0.5rem", padding: "0.5rem" }}
      />

      <textarea
        rows={3}
        value={text}
        onChange={(e) => setText(e.target.value)}
        placeholder="What's happening?"
        style={{ width: "100%", padding: "0.5rem" }}
      />
      <br />
      <button onClick={handlePost} style={{ marginTop: "0.5rem" }}>
        Tweet
      </button>

      <hr />

      {tweets.map(tweet => (
        <div key={tweet.id} style={{ padding: "1rem 0", borderBottom: "1px solid #ccc" }}>
          <p><strong>@{tweet.username}</strong> <em style={{ fontSize: "0.8rem" }}>{new Date(tweet.created_at).toLocaleString()}</em></p>

          {editingId === tweet.id ? (
            <>
              <textarea
                rows={2}
                value={editingText}
                onChange={e => setEditingText(e.target.value)}
                style={{ width: "100%" }}
              />
              <button onClick={handleEditSave}>Save</button>
              <button onClick={() => setEditingId(null)} style={{ marginLeft: "0.5rem" }}>Cancel</button>
            </>
          ) : (
            <>
              <p>{tweet.text}</p>
              <button onClick={() => startEditing(tweet.id, tweet.text)}>Edit</button>
              <button onClick={() => handleDelete(tweet.id)} style={{ marginLeft: "0.5rem" }}>Delete</button>
            </>
          )}
        </div>
      ))}
    </div>
  );
}

export default App;
