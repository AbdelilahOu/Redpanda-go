-- Customers Table
CREATE TABLE Customers (
	CustomerID INT PRIMARY KEY AUTO_INCREMENT,
	FirstName VARCHAR(50) NOT NULL,
	LastName VARCHAR(50) NOT NULL,
	Email VARCHAR(100) UNIQUE,
	Phone VARCHAR(20),
	Address VARCHAR(255),
	City VARCHAR(50),
	State VARCHAR(2),
	ZipCode VARCHAR(10)
);

-- Orders Table
CREATE TABLE Orders (
	OrderID INT PRIMARY KEY AUTO_INCREMENT,
	CustomerID INT,
	OrderDate TIMESTAMP,
	TotalAmount DECIMAL(10, 2),
	Status VARCHAR(20), -- e.g., 'Pending', 'Shipped', 'Delivered'
	FOREIGN KEY (CustomerID) REFERENCES Customers (CustomerID)
);

-- Products Table
CREATE TABLE Products (ProductID INT PRIMARY KEY AUTO_INCREMENT, ProductName VARCHAR(100) NOT NULL, Description TEXT, Price DECIMAL(10, 2), Category VARCHAR(50));

-- Order Items (Line Items)
CREATE TABLE OrderItems (
	OrderItemID INT PRIMARY KEY AUTO_INCREMENT,
	OrderID INT,
	ProductID INT,
	Quantity INT,
	Price DECIMAL(10, 2), -- Price at the time of the order
	FOREIGN KEY (OrderID) REFERENCES Orders (OrderID),
	FOREIGN KEY (ProductID) REFERENCES Products (ProductID)
);

-- Sample Data
INSERT INTO
	Customers (FirstName, LastName, Email, Phone, Address, City, State, ZipCode)
VALUES
	('John', 'Doe', 'john.doe@example.com', '555-123-4567', '123 Main St', 'Anytown', 'CA', '91234'),
	('Jane', 'Smith', 'jane.smith@example.com', '555-987-6543', '456 Oak Ave', 'Springfield', 'IL', '62704'),
	('Robert', 'Jones', 'robert.jones@example.com', '555-555-5555', '789 Pine Ln', 'Hill Valley', 'CA', '90210');

INSERT INTO
	Orders (CustomerID, OrderDate, TotalAmount, Status)
VALUES
	(1, '2023-10-26', 125.50, 'Shipped'),
	(2, '2023-10-27', 75.00, 'Pending'),
	(1, '2023-10-28', 200.00, 'Delivered');

INSERT INTO
	Products (ProductName, Description, Price, Category)
VALUES
	('Widget A', 'A standard widget', 25.00, 'Widgets'),
	('Gadget B', 'A fancy gadget', 50.00, 'Gadgets'),
	('Thingamajig C', 'A useful thingamajig', 100.00, 'Thingamajigs');

INSERT INTO
	OrderItems (OrderID, ProductID, Quantity, Price)
VALUES
	(1, 1, 2, 25.00),
	(1, 2, 1, 50.00),
	(2, 1, 3, 25.00),
	(3, 3, 2, 100.00);
