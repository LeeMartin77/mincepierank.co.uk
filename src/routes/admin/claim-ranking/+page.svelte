<h3>Claim Anonymous pies</h3>
<form method="POST" action="?/claim">
  <div>
    <label for="userid">userid</label>
    <input name="userid" />
  </div>
  <div>
    <label for="anonuserid">anonuserid</label>
    <input name="anonuserid" />
  </div>

  <button>Claim</button>
</form>

<style>
  form {
    display: flex;
    flex-direction: column;
    gap: 0.5em;
  }
  form div {
    display: flex;
  }
  form label {
    margin: 0;
  }
  form input {
    border: 1px solid black;
  }
  button {
    padding: 1em;
  }
</style>
