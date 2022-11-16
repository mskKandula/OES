CREATE TABLE `Examiners`(  
    id INT NOT NULL AUTO_INCREMENT,  
    name VARCHAR(45) NOT NULL,
    age INT,
    email VARCHAR(45) NOT NULL,
    mobileNo VARCHAR(45) NOT NULL,
    password VARCHAR(100) NOT NULL,
    clientId VARCHAR(100) NOT NULL,
    PRIMARY KEY (`id`)  
); 

CREATE TABLE `Students` (  
    id INT NOT NULL AUTO_INCREMENT,  
    name VARCHAR(45) NOT NULL,
    email VARCHAR(45) NOT NULL UNIQUE,
    mobileNo VARCHAR(45) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    clientId VARCHAR(100) NOT NULL,
    PRIMARY KEY (`id`)  
); 

CREATE TABLE `Menu` (  
    id INT NOT NULL AUTO_INCREMENT,  
    name VARCHAR(45) NOT NULL,
    url VARCHAR(45) NOT NULL,
    description VARCHAR(45) NOT NULL,
    PRIMARY KEY (`id`)  
); 

INSERT INTO `Menu`(`name`,`url`,`description`) VALUES
('Dashboard','/dashboard','Dashboard'),
('MultipleStudentsRegistration','/multipleStudentsRegistration','MultipleStudentsRegistration'),
('StudentsList','/studentsList','StudentsList'),
('Uploadquestions','/uploadQuestions','Uploadquestions'),
('UploadVideo','/uploadVideo','UploadVideo'),
('StudentDashboard','/studentDashboard','StudentDashboard'),('Exam','/onlineExam','Exam'),
('VideoContent','/fetchVideos','VideoContent'),('WhiteBoard','/whiteBoard','WhiteBoard'),
('BroadcastVideo','/broadcastVideo?publish=true','BroadcastVideo'),
('BroadcastVideo','/broadcastVideo','BroadcastVideo');

CREATE TABLE `Role` (  
    id INT NOT NULL AUTO_INCREMENT,  
    name VARCHAR(45) NOT NULL,
    PRIMARY KEY (`id`)  
);  

INSERT INTO `Role`(`name`) VALUES('Examiner'),('Student');

CREATE TABLE `UserRole` (  
    userId INT NOT NULL,  
    roleId INT NOT NULL
);

CREATE TABLE `RoleMenu` (  
    roleId INT NOT NULL,  
    menuId INT NOT NULL
);

INSERT INTO `RoleMenu`(`roleId`,`menuId`) VALUES(2,6),(2,7),(2,8),(2,9),(2,11);

INSERT INTO `RoleMenu`(`roleId`,`menuId`) VALUES(1,1),(1,2),(1,3),(1,4),(1,5),(1,9),(2,10);

CREATE TABLE `VideoContent`(
    id INT NOT NULL AUTO_INCREMENT,  
    name VARCHAR(45) NOT NULL,
    videoUrl VARCHAR(100) NOT NULL,
    thumbnailPath VARCHAR(100) NOT NULL,
    contentType VARCHAR(50) NOT NULL,
    description VARCHAR(100),
    clientId VARCHAR(100) NOT NULL,
    PRIMARY KEY (`id`));