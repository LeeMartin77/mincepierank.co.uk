<script lang="ts">
  import LinkButton from '$components/generic/LinkButton.svelte';
</script>

<svelte:head>
  <title>About :: Mince Pie Rank</title>
  <meta name="description" content="All about what Mince Pie Rank is" />
</svelte:head>

<h1>About</h1>
<h3>What is this?</h3>
<p>
  This started as a fun project between me and my partner - where we wanted to rank the mince pies
  we ate over christmas. What was interesting was the second we mentioned the idea to other people,
  they were interested too.
</p>
<h3>Cookie and Privacy Policy Summary</h3>
<p>
  We use a simple authentication cookie once a user is logged in, and only store a user identifier
  (usually an email) so that users know what pies they have ranked. We store no other user data than
  this, and I&apos;ve got no intention of storing more.
</p>
<p>
  <LinkButton href="/about/cookies">Full Cookie Policy</LinkButton>
</p>
<p>
  <LinkButton href="/about/privacy">Full Privacy Policy</LinkButton>
</p>
<p>
  <LinkButton href="https://github.com/LeeMartin77/mincepierank.co.uk">
    Source code is publicly available on Github.
  </LinkButton>
</p>
<h4>Issues? Questions?</h4>
<p>
  <LinkButton href="mailto:piemaster@mincepierank.co.uk">piemaster@mincepierank.co.uk</LinkButton>
</p>
