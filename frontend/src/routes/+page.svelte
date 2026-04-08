<script>
	const PORT = 8080;
	const URL = `http://localhost:${PORT}`;
	const handleClick = async () => {
		const response = await fetch(`${URL}`);
		const data = await response.text();
		console.log(data);
	};
</script>

<button onclick={handleClick}>button</button>

<style>
</style>
