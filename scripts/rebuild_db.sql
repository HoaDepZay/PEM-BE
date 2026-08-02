-- Drop tables if they exist
IF OBJECT_ID('dbo.Expenses', 'U') IS NOT NULL DROP TABLE dbo.Expenses;
IF OBJECT_ID('dbo.Categories', 'U') IS NOT NULL DROP TABLE dbo.Categories;
IF OBJECT_ID('dbo.CategoryGroups', 'U') IS NOT NULL DROP TABLE dbo.CategoryGroups;

-- 1. Create CategoryGroups table
CREATE TABLE CategoryGroups (
    GroupID UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWSEQUENTIALID(),
    UserID UNIQUEIDENTIFIER NULL,
    Name NVARCHAR(100) NOT NULL,
    Icon VARCHAR(50) NOT NULL,
    Color VARCHAR(20) NOT NULL,
    TotalBudget BIGINT NOT NULL DEFAULT 0,
    CreatedAt DATETIME NOT NULL DEFAULT GETDATE(),
    FOREIGN KEY (UserID) REFERENCES Users(UserID) ON DELETE NO ACTION
);

-- 2. Create Categories table
CREATE TABLE Categories (
    CategoryID UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWSEQUENTIALID(),
    GroupID UNIQUEIDENTIFIER NOT NULL,
    UserID UNIQUEIDENTIFIER NULL,
    Name NVARCHAR(100) NOT NULL,
    BudgetType VARCHAR(20) NOT NULL, -- 'DAILY', 'WEEKLY', 'MONTHLY', 'YEARLY'
    BudgetAmount BIGINT NOT NULL DEFAULT 0,
    CreatedAt DATETIME NOT NULL DEFAULT GETDATE(),
    FOREIGN KEY (GroupID) REFERENCES CategoryGroups(GroupID) ON DELETE CASCADE,
    FOREIGN KEY (UserID) REFERENCES Users(UserID) ON DELETE NO ACTION
);

-- 3. Create Expenses table
CREATE TABLE Expenses (
    ExpenseID UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWSEQUENTIALID(),
    UserID UNIQUEIDENTIFIER NOT NULL,
    CategoryID UNIQUEIDENTIFIER NOT NULL,
    Amount BIGINT NOT NULL,
    Note NVARCHAR(500) NULL,
    ImageURL VARCHAR(255) NULL,
    ExpenseDate DATETIME NOT NULL DEFAULT GETDATE(),
    CreatedAt DATETIME NOT NULL DEFAULT GETDATE(),
    FOREIGN KEY (UserID) REFERENCES Users(UserID) ON DELETE CASCADE,
    FOREIGN KEY (CategoryID) REFERENCES Categories(CategoryID) ON DELETE CASCADE
);

-- 4. Create Trigger to update CategoryGroups.TotalBudget
GO
CREATE TRIGGER trg_UpdateTotalBudget
ON Categories
AFTER INSERT, UPDATE, DELETE
AS
BEGIN
    SET NOCOUNT ON;

    -- Update CategoryGroups whose Categories were affected
    UPDATE cg
    SET cg.TotalBudget = ISNULL((SELECT SUM(c.BudgetAmount) FROM Categories c WHERE c.GroupID = cg.GroupID), 0)
    FROM CategoryGroups cg
    WHERE cg.GroupID IN (
        SELECT GroupID FROM inserted
        UNION
        SELECT GroupID FROM deleted
    );
END;
GO
