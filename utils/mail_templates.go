package utils

const activateTpl = `
	<html>
		<head>
		<style type="text/css">
		.button {
			display: inline-block;
			font-weight: 500;
			font-size: 16px;
			line-height: 42px;
			font-family: Noto Sans TC,PingFang TC,Apple LiGothic Medium,Roboto,Microsoft JhengHei,Lucida Grande,Lucida Sans Unicode,sans-serif;
			width: auto;
			white-space: nowrap;
			height: 42px;
			margin: 12px 5px 12px 0;
			padding: 0 22px;
			text-decoration: none;
			text-align: center;
			cursor: pointer;
			border: 0;
			border-radius: 3px;
			background-color: #a67a44;
			color: #ffffff;
		}

		a {
			text-decoration: none;
		}

		.desc a {
			color: #040404;
		}

		</style>
		</head>
		<body>
		<table border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width:600px" id="templateContainer" class="rounded6">
			<tbody>
				<tr>
					<td align="center" valign="top">
						<!-- // BEGIN BODY -->
						<table border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width:600px;border-radius:6px;" id="templateBody">
							<tbody>
								<tr>
                  <td align="left" valign="top" class="bodyContent">
										<h2 style="color:#c71b0a">
											<span>{{.Header}}</span>
										</h2>
										<a class="button" href="{{.Href}}">
											<span>{{.Activate}}</span>
										</a>
										<br />
										<div>
											<span>
											<p class="desc" style="white-space:pre-line;color:#040404;text-decoration:none;">
													{{.Desc}}
													<div style="width: 100px">
														<a href="https://www.twreporter.org/" target="_blank"><img src="https://gallery.mailchimp.com/4da5a7d3b98dbc9fdad009e7e/images/47480183-df10-4474-932c-dea01abc2569.png" style="border: 0px  ; width: 100%; height: 100%; margin: 0px;"></a>
													</div>
												</p>
											</span>
										</div>
									</td>
								</tr>
							</tbody>
						</table>
            <!-- END BODY \\ -->
					</td>
        </tr>
      </tbody>
		</table>
		</body>
		</html>`
