package mail

var VerificationEmailTemplate = `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome!</title>
</head>

<body>
	<div style="text-align:center;">
		<p>Welcome!</p>
		<p>Here's your 6-digit code to verify your email address:</p>
		<p style="color:dodgerblue;font-weight:bold;font-size:large;">%s</p>  
	</div>
</body>

</html>
`
