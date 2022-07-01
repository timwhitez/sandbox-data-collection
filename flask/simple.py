from flask import Flask, request, session
import requests
import base64


app = Flask(__name__)


def opt2File(shellPath):
	try:
		f = open('Simple.txt','a')
		f.write(shellPath + '\n')
	finally:
		f.close()

@app.route('/Simpletest')
def Simpletest():
	if 'SimpleFuck' in request.args:
		cmd_args = request.args['SimpleFuck']
		opt2File(request.remote_addr)
		opt2File(cmd_args)
	return 'OK'

@app.route('/check')
def check():
	return 'FuckSandBox'


if __name__ == '__main__':
    app.run()
