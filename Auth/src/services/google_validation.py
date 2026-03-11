from google.oauth2 import id_token
from google.auth.transport import requests

class GoogleOAUTH:
	def __init__(self, client_id:str) -> None:
		self.google_client_id = client_id

	def verify_google_token(self, token: str):
		try:
			id_info = id_token.verify_oauth2_token(
				token,
				requests.Request(),
				self.google_client_id
			)
			return id_info
		except ValueError:
			return None
		
