CREATE TABLE reports(
    bun.BaseModel 
	ReportID      bigserial PRIMARY KEY       
	UserInfoID    int       
	UserInfo      *UserInfo 
	IssueID       TEXT    
	SpentOn       TEXT    
	Hours         TEXT    
	ActivityID    TEXT    
	Comments      TEXT    
);
CREATE tABLE UserModel(
un.BaseModel 
	MattermostID  TEXT primary key
	RedmineId     INT    
	RedmineApiKey TEXT 

);
CREATE TABLE UserInfo (
    bun.BaseModel 

	UserInfoID INT 

	UserID TEXT 

	CreationDate TIMESTAMP 

	FinishDate TIMESTAMP
);