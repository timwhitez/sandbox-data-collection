def opt2File(shellPath):
	try:
		f = open('result.csv','a')
		f.write(shellPath + '\n')
	finally:
		f.close()


result=[]
mac = ""
pwd0 = ""
ds = ""
hn = ""
pn = ""
fn = ""
un = ""
ip = ""

fi ="ip, Username, Hostname, Filename, Path, Processnum, Disksize"
opt2File(fi)
file = open("1.txt",encoding = "utf-8")
for text in file.readlines():
	data1 = text.strip('\n')
	if data1 != "":
		data1 = data1.split(": ")
		if len(data1)==1:
			ip = data1[0]
		elif data1[0]=="Username":
			un = data1[1]
		elif data1[0] == "Filename":
			fn = data1[1]
		elif data1[0]=="processnum":
			pn = data1[1]
		elif data1[0]=="Hostname":
			hn = data1[1]
		elif data1[0]=="Disksize":
			ds = data1[1]
		elif data1[0]=="Pwd":
			pwd0 = data1[1]

	else:
		result = ip + ","+un.split("\\")[-1] + ","+hn + ","+fn + ","+pwd0 + ","+pn + ","+ds
		opt2File(result)
		pwd0 = ""
		ds = ""
		hn = ""
		pn = ""
		fn = ""
		un = ""
		ip = ""