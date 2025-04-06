/*

---------------------------------------------------------------------------------------------
------------------------------------------USER MANAGEMENT------------------------------------
---------------------------------------------------------------------------------------------

STORED PROCEDURE NAME: CreateUser
STORED PROCEDURE DESCRIPTION: 
    Criação de usário no sistema. Valida todos os campos recebidos pela aplicação se são validos e insere o novo falaor, caso haja falha, nada é executado, e manda a resposta para a aplicação em formato JSON com os detalhes, no attributo de retorno (p_Message)

---------------------------------------------------------------------------------------------
------------------------------------------INCOME MANAGEMENT----------------------------------
---------------------------------------------------------------------------------------------

STORED PROCEDURE NAME: CreateUserParentIncome
TORED PROCEDURE DESCRIPTION: 
    Procedure para criar um novo income para o usuario. Associado na table UserIncomeItem o novo record e criando (baseado na recurency escolhida) todos os records the forecast associados.

STORED PROCEDURE NAME: UpdateUserParentIncome
STORED PROCEDURE DESCRIPTION: 
   Procedure altera valores de income baseados na data (utilizados para modificar um valor de um income a partir de uma data para frente no forecast)
   Tambe utilizada para inativar (income não é mais valido, porém mantem o historico) pela Flag IsActive

STORED PROCEDURE NAME: DeleteUserParentIncomecls
STORED PROCEDURE DESCRIPTION: 
   Deleta um User Parent Income, e todos seus childs associados.

---------------------------------------------------------------------------------------------
------------------------------------------INCOME CHILD---------------------------------------
---------------------------------------------------------------------------------------------
   
STORED PROCEDURE NAME: CreateUserChildIncomeTax
STORED PROCEDURE VERSION: 1.0
STORED PROCEDURE DESCRIPTION: 
   Cria um child income tax. -- NOTA: PRECISA DE IMPROVEMENTS PARA REFLETIR MESMA LOGICA DO CHILD ASSET
   
STORED PROCEDURE NAME: CreateUserChildIncomeExpense
STORED PROCEDURE DESCRIPTION: 
   Cria um child income expense. -- NOTA: PRECISA DE IMPROVEMENTS PARA REFLETIR MESMA LOGICA DO CHILD ASSET

---------------------------------------------------------------------------------------------
-----------------------------------------ASSET MANAGEMENT------------------------------------
---------------------------------------------------------------------------------------------

STORED PROCEDURE NAME: CreateUserAsset
STORED PROCEDURE DESCRIPTION: 
    Procedure para criar uma novo asset associado a um usuario. É obrigatório a seleção de um "AssetType" valido para criação.

STORED PROCEDURE NAME: CreateUserAssetParentIncome
STORED PROCEDURE DESCRIPTION: 
    Procedure para criar um novo income associado a un Asset que, consequentemente, é associado a umo usuario. Criando um novo  UserAssetIncomeItem entity na tabela FinancialUserItem e também criando (baseado na recorrencia escolhida) todos os records the forecast associados na userfinancialforecast.

STORED PROCEDURE NAME: DeleteUserAssetParentIncome
STORED PROCEDURE DESCRIPTION: 
   Deletar o UserAsset. Essa procedure é para um hard delete, excluindo todo historico tanto de forecasts e de actuals (soft deletes são executados pelo UpdateUserAssetParentIncome, pela flag IsActive). Aqui todos os dados são deletados sem possibilidade de recuperação.
   Todos os items relacionado a user Parent serão deletados.
---------------------------------------------------------------------------------------------
------------------------------------------ASSET CHILD MANAGEMENT-----------------------------
---------------------------------------------------------------------------------------------

STORED PROCEDURE NAME: CreateUserAssetChildIncomeTax
STORED PROCEDURE DESCRIPTION: 
   Criar um Tax associado a income existente. A procuedure irá criar um novo record na table FinancialuserItem e replicate (baseado no income existent), todos os records, associado um tax com as mesma datas.

STORED PROCEDURE NAME: CreateUserAssetChildIncomeExpense
STORED PROCEDURE DESCRIPTION: 
   Criar um Expense associado a user asset income existente. A procuedure irá criar um novo record na table FinancialuserItem e replicar (baseado no income existent), todos os records, associado o novo expense com as mesma datas.

STORED PROCEDURE NAME: DeleteUserAssetChildIncomeExpense
STORED PROCEDURE DESCRIPTION: 
     Deletar um expense child User Asset Income. Essa procedure é para um hard delete, excluindo todo historico tanto de forecasts e de actuals. Não existe opção de soft delete em childs Aqui todos os dados são deletados sem possibilidade de recuperação.

STORED PROCEDURE NAME: DeleteUserAssetChildIncomeTax
STORED PROCEDURE DESCRIPTION: 
     Deletar um tax child User Asset Income. Essa procedure é para um hard delete, excluindo todo historico tanto de forecasts e de actuals. Não existe opção de soft delete em childs Aqui todos os dados são deletados sem possibilidade de recuperação.
---------------------------------------------------------------------------------------------
------------------------------------------USER PARENT EXPENSE MANAGEMENT----------------------
---------------------------------------------------------------------------------------------


STORED PROCEDURE NAME: CreateUserParentExpense
STORED PROCEDURE DESCRIPTION: 
    Procedure para criar um novo parent expense para o usuario. Associado na tabela financialuseritem o novo record e criando (baseado na recurency escolhida) todos os records the forecast associados.
*/
  
CREATE OR REPLACE PROCEDURE CreateUser(
    IN p_FirstName VARCHAR(100),
    IN p_LastName VARCHAR(255),
	IN p_EmailAddress VARCHAR(255),
	IN p_UserPassword VARCHAR(150),
    IN p_DateOfBirth DATE,
    OUT p_Message TEXT
)

/* ----------------------------------------------------------------------
STORED PROCEDURE NAME: CreateUser
STORED PROCEDURE VERSION: 1.0
STORED PROCEDURE LAST UPDATED DATE: 30-Mar-2025
STORED PROCEDURE DESCRIPTION: 
    Criação de usário no sistema. Valida todos os campos recebidos pela aplicação se são validos e insere o novo falaor, caso haja falha, nada é executado, e manda a resposta para a aplicação em formato JSON com os detalhes, no attributo de retorno (p_Message)

STORED PROCEDURE TEST CASE:

CALL ('Pedro','Pinto','pfbpinto@hotmai.com','$2a$10$"$2a$10$gFxlAuEUcogaJS4R8jiOiudfVjM3q0H.w9GsLvHPuoIcqKtoH2vE6','2025-07-30')

BACKEND VISUALIZATION:
    Select where EmailAddress='pfbpinto@hotmai.com' 

USER INTERFACE:
Login usando email + password (non-hashed password=123)
----------------------------------------------------*/
LANGUAGE plpgsql
AS $$
DECLARE
    existing_user INT;
BEGIN
    -- Verifica se algum campo obrigatório está nulo ou vazio
    IF TRIM(p_FirstName) = '' OR TRIM(p_LastName) = '' OR p_DateOfBirth IS NULL OR 
       TRIM(p_UserPassword) = '' OR TRIM(p_EmailAddress) = '' IS NULL THEN
        p_Message := '{"status": "fail", "message": "Todos os campos são obrigatórios."}';
        RETURN;
    END IF;

    -- Verifica se o e-mail já existe
    SELECT COUNT(*) INTO existing_user FROM UserProfile WHERE EmailAddress = p_EmailAddress;
    IF existing_user > 0 THEN
        p_Message := '{"status": "fail", "message": "O e-mail já está cadastrado."}';
        RETURN;
    END IF;

    -- Insere o novo usuário
    INSERT INTO UserProfile (FirstName, LastName, DateOfBirth, UserPassword, EmailAddress, UserSubscription)
    VALUES (p_FirstName, p_LastName, p_DateOfBirth, p_UserPassword, p_EmailAddress, '0');

    -- Retorna mensagem de sucesso
    p_Message := '{"status": "success", "message": "Usuário criado com sucesso."}';
END;
$$;




/* INCOME MANAGEMENT PROCEDURES*/

-- STORED PROCEDURE TO CREATE NEW PARENT INCOME
CREATE OR REPLACE PROCEDURE CreateUserParentIncome(
    IN p_UserID INT,
    IN p_FinancialUserItemName VARCHAR(255),
    IN p_RecurrencyID INT,
    IN p_FinancialUserEntityItemID INT,
    IN p_ParentIncomeAmount NUMERIC(15,2),
    IN p_BeginDate DATE,
    OUT p_Message TEXT
)

/* ----------------------------------------------------------------------
STORED PROCEDURE NAME: CreateUserParentIncome
STORED PROCEDURE VERSION: 1.0
STORED PROCEDURE LAST UPDATED DATE: 30-Mar-2025
STORED PROCEDURE DESCRIPTION: 
    Procedure para criar um novo income para o usuario. Associado na table UserIncomeItem o novo record e criando (baseado na recurency escolhida) todos os records the forecast associados.
STORED PROCEDURE TEST CASE(S):

CALL CreateUserParentIncome (13,'Test Income One Time-01',1,5,500,'05-03-2025','')  -- Para Criação de "One Time" incomes
CALL CreateUserParentIncome (13,'Test Income Monthly-01',2,1,10000,'05-03-2025','') -- Para Criação de "Monthly" incomes
CALL CreateUserParentIncome (13,'Test Income Quarterly-01',3,6,15000,'05-03-2025','')-- Para Criação de "Quarterly" incomes
CALL CreateUserParentIncome (13,'Test Income Yearly-01',4,4,25000,'05-03-2025','') -- Para Criação de "Time Yearly" incomes

BACKEND VISUALIZATION:
  
select fo.userfinancialforecastid,F.financialuseritemid,U.UserProfileId,F.userentityid,U.FirstName,U.LastName,F.financialuseritemname,E.EntityName,E.entitytype,FO.userfinancialforecastamount,FO.userfinancialforecastbegindate,FO.userfinancialforecastenddate from UserProfile U
Join financialuseritem F on F.userentityid=U.UserProfileID
right join userfinancialforecast FO on FO.financialuseritemid=F.financialuseritemid
join entity E ON E.entityid=F.entityid
where F.userentityid=13  -- Selecione o UserProfileID do usuario que foi criado os incomes
-- and f.financialuseritemid=59 -- Se quiser granularidade no item que foi criado, use o financialuseritemid associado ao Income
order by 1,10

USER INTERFACE:
Navegue até Income, e veja os records no sitema
----------------------------------------------------*/
LANGUAGE plpgsql
AS $$
DECLARE 
    v_NewFinancialUserItemid INT;
    v_CurrentDate DATE;
    v_NextDate DATE;
    v_Iterations INT;
    v_Increment INTERVAL;
    i INT;
BEGIN
    -- Validações de campos obrigatórios
    IF p_UserID IS NULL THEN 
        p_Message := '{"status": "fail", "message": "Missing UserID"}';
        RETURN;
    END IF;
    IF p_FinancialUserItemName IS NULL OR p_FinancialUserItemName = '' THEN 
        p_Message := '{"status": "fail", "message": "Missing FinancialUserItemName"}';
        RETURN;
    END IF;
    IF p_RecurrencyID IS NULL THEN 
        p_Message := '{"status": "fail", "message": "Missing RecurrencyID"}';
        RETURN;
    END IF;
    IF p_FinancialUserEntityItemID IS NULL THEN 
        p_Message := '{"status": "fail", "message": "Missing FinancialUserEntityItemID"}';
        RETURN;
    END IF;
    IF p_ParentIncomeAmount IS NULL OR p_ParentIncomeAmount <= 0 THEN 
        p_Message := '{"status": "fail", "message": "Invalid ParentIncomeAmount"}';
        RETURN;
    END IF;
    IF p_BeginDate IS NULL THEN 
        p_Message := '{"status": "fail", "message": "Missing BeginDate"}';
        RETURN;
    END IF;

    -- Inicia transação manualmente
    BEGIN
        -- Insere um novo Parent Income
    INSERT INTO FinancialUserItem (
        FinancialUserItemName, EntityID, UserEntityID, RecurrencyID, 
        FinancialUserEntityItemID, ParentFinancialUserItemID
    ) VALUES (
        p_FinancialUserItemName, 5, p_UserID, p_RecurrencyID, 
        p_FinancialUserEntityItemID, NULL
    ) RETURNING FinancialUserItemID INTO v_NewFinancialUserItemid;

        -- Configura número de inserções e intervalo conforme a recorrência
        IF p_RecurrencyID = 1 THEN -- One time
            v_Iterations := 1;
            v_Increment := '1 day'::INTERVAL;
        ELSIF p_RecurrencyID = 2 THEN -- Monthly
            v_Iterations := 12;
            v_Increment := '1 month'::INTERVAL;
        ELSIF p_RecurrencyID = 3 THEN -- Quarterly
            v_Iterations := 4;
            v_Increment := '4 months'::INTERVAL;
        ELSIF p_RecurrencyID = 4 THEN -- Yearly
            v_Iterations := 1;
            v_Increment := '1 year'::INTERVAL;
        ELSE
            RAISE EXCEPTION 'Invalid RecurrencyID';
        END IF;

        -- Insere múltiplos registros conforme a recorrência
        v_CurrentDate := p_BeginDate;
        FOR i IN 1..v_Iterations LOOP
            -- Define a data do próximo registro para calcular a end date
            v_NextDate := v_CurrentDate + v_Increment;
            
            INSERT INTO UserFinancialForecast (
                usercategoryid, financialuseritemid, userfinancialforecastbegindate, 
                userfinancialforecastenddate, userfinancialforecastamount, currencyid
            ) VALUES (
                NULL, v_NewFinancialUserItemid, v_CurrentDate,
                CASE 
                    WHEN p_RecurrencyID = 1 THEN v_CurrentDate + INTERVAL '1 day' - INTERVAL '1 day' -- One Time
                    WHEN p_RecurrencyID = 2 AND i < v_Iterations THEN v_NextDate - INTERVAL '1 day' -- Monthly
                    WHEN p_RecurrencyID = 2 AND i = v_Iterations THEN (v_CurrentDate + INTERVAL '1 month') - INTERVAL '1 day' -- Último mês
                    WHEN p_RecurrencyID = 3 AND i < v_Iterations THEN v_NextDate - INTERVAL '1 day' -- Quarterly
                    WHEN p_RecurrencyID = 3 AND i = v_Iterations THEN (v_CurrentDate + INTERVAL '4 months') - INTERVAL '1 day' -- Último trimestre
                    WHEN p_RecurrencyID = 4 THEN v_CurrentDate + INTERVAL '1 year' - INTERVAL '1 day' -- Yearly
                END,
                p_ParentIncomeAmount, 1
            );
            
            -- Atualiza data para a próxima recorrência
            v_CurrentDate := v_NextDate;
        END LOOP;

        -- Se tudo deu certo, define mensagem de sucesso
        p_Message := '{"status": "success", "message": "New user income and forecast created successfully."}';
    EXCEPTION 
        WHEN OTHERS THEN
            -- Captura erro e define mensagem de falha
            p_Message := '{"status": "fail", "message": "Error creating forecast values: ' || SQLERRM || '"}';
    END;
END;
$$;

-- STORED PROCEDURE TO UPDATE PARENT INCOME
-- INCLUIDO INATIVAÇÃO DO INCOME (FLAG ACTIVE OU FALSE)

CREATE OR REPLACE PROCEDURE UpdateUserParentIncome(
    IN p_FinancialUserItemID INT,
    IN p_UserID INT,
    IN p_NewFinancialUserItemName VARCHAR(255),
    IN p_NewParentIncomeAmount NUMERIC(15,2),
    IN p_NewBeginDate DATE,
    IN p_IsActive BOOLEAN,
    OUT p_Message TEXT
)
LANGUAGE plpgsql
AS $$
DECLARE 
    v_FinancialItemExist INT;
    v_AssociatedActual INT;
BEGIN
    -- Validações de campos obrigatórios
    IF p_FinancialUserItemID IS NULL THEN 
        p_Message := '{"status": "fail", "message": "Missing FinancialUserItemID"}';
        RETURN;
    END IF;
    IF p_UserID IS NULL THEN 
        p_Message := '{"status": "fail", "message": "Missing UserID"}';
        RETURN;
    END IF;
    IF p_NewFinancialUserItemName IS NULL OR p_NewFinancialUserItemName = '' THEN 
        p_Message := '{"status": "fail", "message": "Missing the new FinancialUserItemName"}';
        RETURN;
    END IF;
    IF p_NewParentIncomeAmount IS NULL THEN 
        p_Message := '{"status": "fail", "message": "Missing the new ParentIncomeAmount"}';
        RETURN;
    END IF;
    IF p_NewBeginDate IS NULL THEN 
        p_Message := '{"status": "fail", "message": "Missing new BeginDate"}';
        RETURN;
    END IF;

    -- Verifica se o FinancialUserItem existe
    SELECT financialuseritemid 
    INTO v_FinancialItemExist
    FROM UserFinancialForecast
    WHERE financialuseritemid = p_FinancialUserItemID
    LIMIT 1;

    IF v_FinancialItemExist IS NULL THEN
        p_Message := '{"status": "fail", "message": "UserParentIncome not found"}';
        RETURN;
    END IF;

    -- Confere se o usuário é dono do financial item
    IF p_UserID <> (
        SELECT userentityid 
        FROM FinancialUserItem 
        WHERE FinancialUserItemID = p_FinancialUserItemID
    ) THEN
        p_Message := '{"status": "fail", "message": "User does not match FinancialUserItem"}';
        RETURN;
    END IF;

    IF p_IsActive = TRUE THEN
        -- Atualiza nome
        IF p_NewFinancialUserItemName IS NOT NULL AND p_NewFinancialUserItemName <> '' THEN
            UPDATE FinancialUserItem
            SET FinancialUserItemName = p_NewFinancialUserItemName
            WHERE FinancialUserItemID = p_FinancialUserItemID;
        END IF;

        -- Atualiza os valores das previsões futuras
        IF p_NewBeginDate IS NOT NULL AND p_NewParentIncomeAmount IS NOT NULL THEN
            UPDATE UserFinancialForecast
            SET userfinancialforecastamount = p_NewParentIncomeAmount
            WHERE financialuseritemid = p_FinancialUserItemID
            AND userfinancialforecastbegindate >= p_NewBeginDate;
        END IF;

        p_Message := '{"status": "success", "message": "UserParentIncome updated successfully."}';
    END IF;

    IF p_IsActive = FALSE THEN
        -- Verifica se existe Forecast futuro com actual associado
        SELECT fr.userfinancialforecastid INTO v_AssociatedActual
        FROM userforecastactualrelation fr
        WHERE fr.userfinancialforecastid IN (
            SELECT UserFinancialForecastid
            FROM UserFinancialForecast f
            WHERE f.FinancialUserItemID = p_FinancialUserItemID
            AND f.userfinancialforecastbegindate > p_NewBeginDate
        )
        LIMIT 1;

        IF v_AssociatedActual IS NULL THEN    
            -- Inativa o item
            UPDATE FinancialUserItem
            SET IsActive = FALSE
            WHERE financialuseritemid = p_FinancialUserItemID;

            -- Deleta os forecasts futuros
            DELETE FROM UserFinancialForecast
            WHERE financialuseritemid = p_FinancialUserItemID
            AND userfinancialforecastbegindate > p_NewBeginDate;

            p_Message := '{"status": "success", "message": "Forecasting for this UserParentIncome is inactive"}';
        ELSE 
            p_Message := '{"status": "fail", "message": "There is an actual record associated with a forecast in a future date. Delete the actual or adjust the date. Inactivation aborted."}';
        END IF;
    END IF;
END;
$$;


-- DELETE PARENT INCOME

CREATE OR REPLACE PROCEDURE DeleteUserParentIncome(
    IN p_FinancialUserItemID INT,
    IN p_UserID INT,
    OUT p_Message TEXT
)
/* ----------------------------------------------------------------------
STORED PROCEDURE NAME: DeleteUserParentIncome
STORED PROCEDURE VERSION: 1.0
STORED PROCEDURE LAST UPDATED DATE: 30-Mar-2025
STORED PROCEDURE DESCRIPTION: 
   Deleta um User Parent Income, e todos seus childs associados.   -- NOTA: PRECISA DE IMPROVEMENTS PARA REFLETIR MESMA LOGICA DO DELETE USER ASSET
STORED PROCEDURE TEST CASE(S):


BACKEND VISUALIZATION:
  
USER INTERFACE:

----------------------------------------------------*/
LANGUAGE plpgsql
AS $$
BEGIN
    -- Verifica se o FinancialUserItem existe
    IF NOT EXISTS (SELECT 1 FROM FinancialUserItem WHERE FinancialUserItemID = p_FinancialUserItemID AND UserEntityID = p_UserID) THEN
        p_Message := '{"status": "fail", "message": "UserParentIncome not found"}';
        RETURN;
    END IF;

    -- Remove forecasts associadas
    DELETE FROM UserFinancialForecast WHERE FinancialUserItemID = p_FinancialUserItemID;

    --Remove actuals associados

    DELETE FROM UserFinancialActual WHERE FinancialUserItemID = p_FinancialUserItemID;

    -- Remove o FinancialUserItem
    DELETE FROM FinancialUserItem WHERE FinancialUserItemID = p_FinancialUserItemID;

    p_Message := '{"status": "success", "message": "UserParentIncome deleted successfully."}';
END;
$$;


-- STORED PROCEDURE TO CREATE NEW CHIID INCOME TAX PARA UM PARENT INCOME

CREATE OR REPLACE PROCEDURE CreateUserChildIncomeTax(
    IN p_UserID INT,
    IN p_FinancialUserItemName VARCHAR(255),
	IN p_RecurrencyID INT,
	IN p_FinancialUserEntityItemID INT,
	IN p_ParentFinancialUserItemID INT,
    OUT p_Message TEXT
)

/* ----------------------------------------------------------------------
STORED PROCEDURE NAME: CreateUserChildIncomeTax
STORED PROCEDURE VERSION: 1.0
STORED PROCEDURE LAST UPDATED DATE: 30-Mar-2025
STORED PROCEDURE DESCRIPTION: 
   Cria um child income tax. -- NOTA: PRECISA DE IMPROVEMENTS PARA REFLETIR MESMA LOGICA DO CHILD ASSET
STORED PROCEDURE TEST CASE(S):

BACKEND VISUALIZATION:
  
USER INTERFACE:

----------------------------------------------------*/
LANGUAGE plpgsql
AS $$

BEGIN
    -- Verifica se algum campo obrigatório está nulo ou vazio
    IF p_UserID IS NULL THEN
        p_Message := '{"status": "fail", "message": "Missing UserProfileID}';
        RETURN;
    END IF;

    IF p_ParentFinancialUserItemID IS NULL THEN
        p_Message := '{"status": "fail", "message": "Parent Income Item is mandatory}';
        RETURN;
    END IF;
    -- Insere o novo Child Tax Income
    INSERT INTO FinancialUserItem (
        FinancialUserItemName, EntityID, UserEntityID, RecurrencyID, FinancialUserEntityItemID, ParentFinancialUserItemID
    ) VALUES (
        p_FinancialUserItemName, 7, p_UserID, p_RecurrencyID, p_FinancialUserEntityItemID, p_ParentFinancialUserItemID
    );
    -- Retorna mensagem de sucesso
    p_Message := '{"status": "success", "message": "New user tax associated to the parent income created successfully."}';
END;
$$;



-- STORED PROCEDURE TO CREATE NEW CHILD INCOME EXPENSE FOR PARENT INCOME


CREATE OR REPLACE PROCEDURE CreateUserChildIncomeExpense(
    IN p_UserID INT,
    IN p_FinancialUserItemName VARCHAR(255),
	IN p_RecurrencyID INT,
	IN p_FinancialUserEntityItemID INT,
	IN p_ParentFinancialUserItemID INT,
    OUT p_Message TEXT
)

/* ----------------------------------------------------------------------
STORED PROCEDURE NAME: CreateUserChildIncomeExpense
STORED PROCEDURE VERSION: 1.0
STORED PROCEDURE LAST UPDATED DATE: 30-Mar-2025
STORED PROCEDURE DESCRIPTION: 
   Cria um child income expense. -- NOTA: PRECISA DE IMPROVEMENTS PARA REFLETIR MESMA LOGICA DO CHILD ASSET
STORED PROCEDURE TEST CASE(S):

BACKEND VISUALIZATION:
  
USER INTERFACE:

----------------------------------------------------*/
LANGUAGE plpgsql
AS $$

BEGIN
    -- Verifica se algum campo obrigatório está nulo ou vazio
    IF p_UserID IS NULL THEN
        p_Message := '{"status": "fail", "message": "Missing UserProfileID}';
        RETURN;
    END IF;

    IF p_ParentFinancialUserItemID IS NULL THEN
        p_Message := '{"status": "fail", "message": "Parent Income Item is mandatory}';
        RETURN;
    END IF;
    -- Insere o novo Child Expense Income
    INSERT INTO FinancialUserItem (
        FinancialUserItemName, EntityID, UserEntityID, RecurrencyID, FinancialUserEntityItemID, ParentFinancialUserItemID
    ) VALUES (
        p_FinancialUserItemName, 8, p_UserID, p_RecurrencyID, p_FinancialUserEntityItemID, p_ParentFinancialUserItemID
    );

    -- Retorna mensagem de sucesso
    p_Message := '{"status": "success", "message": "New user expense associated to the parent income created successfully."}';
END;
$$;



/* ASSET MANAGEMENT */

-- CREATING A NEW USER ASSET

CREATE OR REPLACE PROCEDURE CreateUserAsset(
    IN p_AssetTypeID INT,
    IN p_UserProfileID INT,
    IN p_UserAssetName VARCHAR(100),
    IN p_UserAssetValueAmount DECIMAL(15,2),
    IN p_UserAssetAcquisitionBeginDate DATE,
    IN p_UserAssetAcquisitionEndDate DATE,
    OUT p_Message TEXT
)

/* ----------------------------------------------------------------------
STORED PROCEDURE NAME: CreateUserAsset
STORED PROCEDURE VERSION: 1.0
STORED PROCEDURE LAST UPDATED DATE: 30-Mar-2025
STORED PROCEDURE DESCRIPTION: 
    Procedure para criar uma novo asset associado a um usuario. É obrigatório a seleção de um "AssetType" valido para criação.
STORED PROCEDURE TEST CASE(S):

CALL CreateUserAsset (1,1,'Novo Apartamento',500000,'2025-01-01',null,'')

BACKEND VISUALIZATION:
  

select * from userasset where UserProfileID=1 -- Mostra os assets criados para o usuario, onde o novo deve estar sendo listado.

USER INTERFACE:

TBD
----------------------------------------------------*/
LANGUAGE plpgsql
AS $$
BEGIN
    -- Validação dos campos obrigatórios
    IF p_AssetTypeID IS NULL THEN
        p_Message := '{"status": "fail", "message": "Missing AssetTypeID"}';
        RETURN;
    END IF;

    IF p_UserProfileID IS NULL THEN
        p_Message := '{"status": "fail", "message": "Missing UserProfileID"}';
        RETURN;
    END IF;

    IF p_UserAssetValueAmount IS NULL OR p_UserAssetValueAmount < 0 THEN
        p_Message := '{"status": "fail", "message": "Missing or invalid UserAssetValueAmount"}';
        RETURN;
    END IF;

    IF p_UserAssetAcquisitionBeginDate IS NULL THEN
        p_Message := '{"status": "fail", "message": "Missing UserAssetAcquisitionBeginDate"}';
        RETURN;
    END IF;

    -- Inserir o novo asset na tabela UserAsset
    INSERT INTO UserAsset (
        AssetTypeID, 
        UserProfileID, 
        UserAssetName, 
        UserAssetValueAmount, 
        UserAssetAcquisitionBeginDate, 
        UserAssetAcquisitionEndDate
    )
    VALUES (
        p_AssetTypeID, 
        p_UserProfileID, 
        p_UserAssetName, 
        p_UserAssetValueAmount, 
        p_UserAssetAcquisitionBeginDate, 
        p_UserAssetAcquisitionEndDate
    );

    p_Message := '{"status": "success", "message": "UserAsset created successfully."}';

END;
$$;


-- CREATE A NEW ASSET PARENT INCOME

CREATE OR REPLACE PROCEDURE CreateUserAssetParentIncome(
    IN p_UserID INT,
    IN p_UserAssetID INT,
    IN p_FinancialUserItemName VARCHAR(255),
    IN p_RecurrencyID INT,
    IN p_FinancialUserEntityItemID INT,
    IN p_ParentIncomeAmount NUMERIC(15,2),
    IN p_BeginDate DATE,
    OUT p_Message TEXT
)
/* ----------------------------------------------------------------------
STORED PROCEDURE NAME: CreateUserAssetParentIncome
STORED PROCEDURE VERSION: 1.0
STORED PROCEDURE LAST UPDATED DATE: 30-Mar-2025
STORED PROCEDURE DESCRIPTION: 
    Procedure para criar um novo income associado a un Asset que, consequentemente, é associado a umo usuario. Criando um novo  UserAssetIncomeItem entity na tabela FinancialUserItem e também criando (baseado na recorrencia escolhida) todos os records the forecast associados na userfinancialforecast.
STORED PROCEDURE TEST CASE(S):

call CreateUserAssetParentIncome (1,2,'Aluguel novo AP',2,2,4500,'2025-03-05','')-- Para Criação de "Monthly" incomes

BACKEND VISUALIZATION:
  
select fo.userfinancialforecastid,F.financialuseritemid,U.UserProfileId,F.userentityid,U.FirstName,U.LastName,F.financialuseritemname,E.EntityName,E.entitytype,FO.userfinancialforecastamount,FO.userfinancialforecastbegindate,FO.userfinancialforecastenddate from UserProfile U
Join financialuseritem F on F.userentityid=U.UserProfileID
right join userfinancialforecast FO on FO.financialuseritemid=F.financialuseritemid
join entity E ON E.entityid=F.entityid
where F.userentityid=2  -- Selecione o UserAssetID que pertence ao usuario que foi criado os incomes
-- and f.financialuseritemid=59 -- Se quiser granularidade no item que foi criado, use o financialuseritemid associado ao Income
order by 1,10

USER INTERFACE:
TBD
----------------------------------------------------*/
LANGUAGE plpgsql
AS $$
DECLARE 
    v_UserValidation INT;
    NewFinancialUserItemid INT;
    v_CurrentDate DATE;
    v_NextDate DATE;
    v_Iterations INT;
    v_Increment INTERVAL;
    i INT;
BEGIN
    -- Validações de campos obrigatórios
    IF p_UserID IS NULL THEN 
        p_Message := '{"status": "fail", "message": "Missing UserID"}';
        RETURN;
    END IF;
    IF p_UserAssetID IS NULL THEN 
        p_Message := '{"status": "fail", "message": "Missing UserAssetID"}';
        RETURN;
    END IF;
    SELECT userprofileid 
    INTO v_UserValidation
    FROM userasset
    WHERE userassetid = p_UserAssetID;
    IF p_UserID <> v_UserValidation THEN
        p_Message := '{"status": "fail", "message": "UserID provided different from the asset owner"}';
        RETURN;
    END IF;
    IF p_FinancialUserItemName IS NULL OR p_FinancialUserItemName = '' THEN 
        p_Message := '{"status": "fail", "message": "Missing FinancialUserItemName"}';
        RETURN;
    END IF;
    IF p_RecurrencyID IS NULL THEN 
        p_Message := '{"status": "fail", "message": "Missing RecurrencyID"}';
        RETURN;
    END IF;
    IF p_FinancialUserEntityItemID IS NULL THEN 
        p_Message := '{"status": "fail", "message": "Missing FinancialUserEntityItemID"}';
        RETURN;
    END IF;
    IF p_ParentIncomeAmount IS NULL OR p_ParentIncomeAmount <= 0 THEN 
        p_Message := '{"status": "fail", "message": "Invalid ParentIncomeAmount"}';
        RETURN;
    END IF;
    IF p_BeginDate IS NULL THEN 
        p_Message := '{"status": "fail", "message": "Missing BeginDate"}';
        RETURN;
    END IF;

    -- Inicia transação manualmente
    BEGIN
        -- Insere um novo Parent Income
        INSERT INTO FinancialUserItem (
            FinancialUserItemName, EntityID, UserEntityID, RecurrencyID, 
            FinancialUserEntityItemID, ParentFinancialUserItemID
        ) VALUES (
            p_FinancialUserItemName, 11, p_UserAssetID, p_RecurrencyID, 
            p_FinancialUserEntityItemID, NULL
        ) RETURNING FinancialUserItemID INTO NewFinancialUserItemid;

        -- Configura número de inserções e intervalo conforme a recorrência
        IF p_RecurrencyID = 1 THEN -- One time
            v_Iterations := 1;
            v_Increment := '1 day'::INTERVAL;
        ELSIF p_RecurrencyID = 2 THEN -- Monthly
            v_Iterations := 12;
            v_Increment := '1 month'::INTERVAL;
        ELSIF p_RecurrencyID = 3 THEN -- Quarterly
            v_Iterations := 4;
            v_Increment := '4 months'::INTERVAL;
        ELSIF p_RecurrencyID = 4 THEN -- Yearly
            v_Iterations := 1;
            v_Increment := '1 year'::INTERVAL;
        ELSE
            RAISE EXCEPTION 'Invalid RecurrencyID';
        END IF;

        -- Insere múltiplos registros conforme a recorrência
        v_CurrentDate := p_BeginDate;
        FOR i IN 1..v_Iterations LOOP
            -- Define a data do próximo registro para calcular a end date
            v_NextDate := v_CurrentDate + v_Increment;
            
            INSERT INTO UserFinancialForecast (
                usercategoryid, financialuseritemid, userfinancialforecastbegindate, 
                userfinancialforecastenddate, userfinancialforecastamount, currencyid
            ) VALUES (
                NULL, NewFinancialUserItemid, v_CurrentDate, 
                CASE 
                    WHEN p_RecurrencyID = 1 THEN v_CurrentDate + INTERVAL '1 day' - INTERVAL '1 day' -- One Time
                    WHEN p_RecurrencyID = 2 AND i < v_Iterations THEN v_NextDate - INTERVAL '1 day' -- Monthly
                    WHEN p_RecurrencyID = 2 AND i = v_Iterations THEN (v_CurrentDate + INTERVAL '1 month') - INTERVAL '1 day' -- Último mês
                    WHEN p_RecurrencyID = 3 AND i < v_Iterations THEN v_NextDate - INTERVAL '1 day' -- Quarterly
                    WHEN p_RecurrencyID = 3 AND i = v_Iterations THEN (v_CurrentDate + INTERVAL '4 months') - INTERVAL '1 day' -- Último trimestre
                    WHEN p_RecurrencyID = 4 THEN v_CurrentDate + INTERVAL '1 year' - INTERVAL '1 day' -- Yearly
                END,
                p_ParentIncomeAmount, 1
            );
            
            -- Atualiza data para a próxima recorrência
            v_CurrentDate := v_NextDate;
        END LOOP;

        -- Se tudo deu certo, define mensagem de sucesso
        p_Message := '{"status": "success", "message": "New user income and forecast created successfully."}';
    EXCEPTION 
        WHEN OTHERS THEN
            -- Captura erro e define mensagem de falha
            p_Message := '{"status": "fail", "message": "Error creating forecast values: ' || SQLERRM || '"}';
    END;
END;
$$;




-- STORED PROCEDURE NAME: DeleteUserAssetParentIncome
-- STORED PROCEDURE VERSION: 1.0
-- LAST UPDATED: 31-Mar-2025
-- DESCRIPTION: 
-- Deleta permanentemente o UserAssetParentIncome (hard delete).
-- Remove todo histórico relacionado (forecasts, actuals, relations).
-- Esse procedimento remove todos os itens relacionados ao Parent Asset.

-- TEST CASE:
-- CALL DeleteUserAssetParentIncome(59, 1, 2, '');

CREATE OR REPLACE PROCEDURE DeleteUserAssetParentIncome(
    IN p_FinancialUserItemID INT,
    IN p_UserID INT,
    IN p_UserAssetID INT, 
    OUT p_Message TEXT
)
LANGUAGE plpgsql
AS $$
DECLARE 
    v_UserValidation INT;
    v_ChildCheck INT;
BEGIN
    -- Validação: FinancialUserItem existe e pertence ao UserAssetID
    IF NOT EXISTS (
        SELECT 1 FROM FinancialUserItem 
        WHERE FinancialUserItemID = p_FinancialUserItemID 
        AND UserEntityID = p_UserAssetID
    ) THEN
        p_Message := '{"status": "fail", "message": "UserAssetParentIncome not found"}';
        RETURN;
    END IF;

    -- Validação: UserAsset pertence ao UserID
    SELECT UserProfileID 
    INTO v_UserValidation
    FROM UserAsset
    WHERE UserAssetID = p_UserAssetID;

    IF p_UserID <> v_UserValidation THEN
        p_Message := '{"status": "fail", "message": "UserID provided is different from the asset owner"}';
        RETURN;
    END IF;

    BEGIN 
        -- Verifica se há Childs relacionados ao Parent
        SELECT parentfinancialuseritemid
        INTO v_ChildCheck
        FROM financialuseritem
        WHERE parentfinancialuseritemid = p_FinancialUserItemID        
        LIMIT 1;
    
        IF v_ChildCheck IS NOT NULL THEN
            -- Tabelas temporárias para armazenar dados relacionados
            CREATE TEMP TABLE TempForecasts AS 
            SELECT f.FinancialUserItemID, uf.UserFinancialForecastID 
            FROM FinancialUserItem f
            JOIN UserFinancialForecast uf ON f.FinancialUserItemID = uf.FinancialUserItemID
            WHERE f.ParentFinancialUserItemID = p_FinancialUserItemID;

            CREATE TEMP TABLE TempActuals AS 
            SELECT f.FinancialUserItemID, ua.UserFinancialActualID 
            FROM FinancialUserItem f
            JOIN UserFinancialActual ua ON f.FinancialUserItemID = ua.FinancialUserItemID
            WHERE f.ParentFinancialUserItemID = p_FinancialUserItemID;

            -- Deletar relações Forecast-Actual
            DELETE FROM UserForecastActualRelation
            WHERE UserFinancialForecastID IN (SELECT UserFinancialForecastID FROM TempForecasts)
               OR UserFinancialActualID IN (SELECT UserFinancialActualID FROM TempActuals);

            -- Deletar Forecasts e Actuals filhos
            DELETE FROM UserFinancialForecast 
            WHERE UserFinancialForecastID IN (SELECT UserFinancialForecastID FROM TempForecasts);

            DELETE FROM UserFinancialActual 
            WHERE UserFinancialActualID IN (SELECT UserFinancialActualID FROM TempActuals);

            -- Deletar FinancialUserItem filhos
            DELETE FROM FinancialUserItem 
            WHERE FinancialUserItemID IN (SELECT FinancialUserItemID FROM TempForecasts)
               OR FinancialUserItemID IN (SELECT FinancialUserItemID FROM TempActuals);
        END IF;

        -- Deletar relações do Parent
        DELETE FROM UserForecastActualRelation
        WHERE UserFinancialForecastID IN (
            SELECT UserFinancialForecastID FROM UserFinancialForecast WHERE FinancialUserItemID = p_FinancialUserItemID
        )
        OR UserFinancialActualID IN (
            SELECT UserFinancialActualID FROM UserFinancialActual WHERE FinancialUserItemID = p_FinancialUserItemID
        );

        -- Deletar Forecasts e Actuals do Parent
        DELETE FROM UserFinancialForecast
        WHERE FinancialUserItemID = p_FinancialUserItemID;

        DELETE FROM UserFinancialActual
        WHERE FinancialUserItemID = p_FinancialUserItemID;

        -- Deletar FinancialUserItem do Parent
        DELETE FROM FinancialUserItem 
        WHERE FinancialUserItemID = p_FinancialUserItemID;

        -- Limpeza
        DROP TABLE IF EXISTS TempForecasts;
        DROP TABLE IF EXISTS TempActuals;

        -- Mensagem de sucesso
        p_Message := '{"status": "success", "message": "User Asset Parent Income and all references deleted successfully."}';

    EXCEPTION 
        WHEN OTHERS THEN 
            p_Message := format('{"status": "fail", "message": "An error occurred: %s"}', SQLERRM);
            DROP TABLE IF EXISTS TempForecasts;
            DROP TABLE IF EXISTS TempActuals;
    END;

END;
$$;



-- Create User Asset Child Income Tax

CREATE OR REPLACE PROCEDURE CreateUserAssetChildIncomeTax(
    IN p_UserID INT,
    IN p_UserAssetID INT,
    IN p_FinancialUserItemName VARCHAR(255),
    IN p_FinancialUserEntityItemID INT,
    IN p_ParentFinancialUserItemID INT,
    IN p_TaxIncomeAmount NUMERIC(15,2),
    OUT p_Message TEXT
)
/* ----------------------------------------------------------------------
STORED PROCEDURE NAME: CreateUserAssetChildIncomeTax
STORED PROCEDURE VERSION: 1.0
STORED PROCEDURE LAST UPDATED DATE: 30-Mar-2025
STORED PROCEDURE DESCRIPTION: 
   Criar um Tax associado a income existente. A procuedure irá criar um novo record na table FinancialuserItem e replicate (baseado no income existent), todos os records, associado um tax com as mesma datas.
STORED PROCEDURE TEST CASE(S):

CALL CreateUserAssetChildIncomeTax(1, 2,'New Child Tax IRPJ',3,60,950,'');

BACKEND VISUALIZATION:
  
select * from financialuseritem where financialuseritemname='New Child Tax IRPJ' -- mostrar o novo record (associado com o seu parent income)
select * from userfinancialforecast where financialuseritemid = (select financialuseritemid from financialuseritem where financialuseritemname='New Child Tax IRPJ') -- Mostrar todos os records na forecast


USER INTERFACE:

TBD
----------------------------------------------------*/
LANGUAGE plpgsql
AS $$

DECLARE 
    v_UserValidation INT;
    v_FinancialUserItemID INT;
BEGIN
    -- Verifica se algum campo obrigatório está nulo ou vazio
    IF p_UserID IS NULL THEN
        p_Message := '{"status": "fail", "message": "Missing UserID"}';
        RETURN;
    END IF;

    IF p_ParentFinancialUserItemID IS NULL THEN
        p_Message := '{"status": "fail", "message": "Parent Income Item is mandatory"}';
        RETURN;
    END IF;

    -- Validação do UserAssetID
    SELECT userprofileid  
    INTO v_UserValidation
    FROM userasset
    WHERE userassetid = p_UserAssetID;

    IF p_UserID <> v_UserValidation THEN
        p_Message := '{"status": "fail", "message": "UserID provided is different from the asset owner"}';
        RETURN;
    END IF;

    -- Insere o novo Child Tax Income
    INSERT INTO FinancialUserItem (
        FinancialUserItemName, EntityID, UserEntityID, RecurrencyID, FinancialUserEntityItemID, ParentFinancialUserItemID
    ) 
    VALUES (
        p_FinancialUserItemName, 12, p_UserAssetID, 1, p_FinancialUserEntityItemID, p_ParentFinancialUserItemID
    )
    RETURNING FinancialUserItemID INTO v_FinancialUserItemID;

    -- Criação dos valores associados ao forecast
    BEGIN
        -- Para cada forecast existente, crie um novo, substituindo os valores conforme necessário
        INSERT INTO userfinancialforecast (
            usercategoryid, financialuseritemid, userfinancialforecastbegindate, 
            userfinancialforecastenddate, userfinancialforecastamount, currencyid
        )
        SELECT 
            usercategoryid, 
            v_FinancialUserItemID,  -- Novo FinancialUserItemID
            userfinancialforecastbegindate, 
            userfinancialforecastenddate, 
            p_TaxIncomeAmount,  -- Valor do imposto
            currencyid
        FROM userfinancialforecast
        WHERE financialuseritemid = p_ParentFinancialUserItemID;

        -- Se não houver erro, envia mensagem de sucesso
        p_Message := '{"status": "success", "message": "New user tax associated to the parent income created successfully."}';
    EXCEPTION
        WHEN OTHERS THEN
            -- Em caso de erro, reverte as mudanças e retorna erro
            ROLLBACK;
            p_Message := '{"status": "fail", "message": "Could not create the Forecast, rolling back changes"}';
    END;

END;
$$;



-- Create User Asset Child Expense Tax

CREATE OR REPLACE PROCEDURE CreateUserAssetChildIncomeExpense(
    IN p_UserID INT,
    IN p_UserAssetID INT,
    IN p_FinancialUserItemName VARCHAR(255),
    IN p_FinancialUserEntityItemID INT,
    IN p_ParentFinancialUserItemID INT,
    IN p_ExpenseAmount NUMERIC(15,2),
    OUT p_Message TEXT
)
/* ----------------------------------------------------------------------
STORED PROCEDURE NAME: CreateUserAssetChildIncomeExpense
STORED PROCEDURE VERSION: 1.0
STORED PROCEDURE LAST UPDATED DATE: 30-Mar-2025
STORED PROCEDURE DESCRIPTION: 
   Criar um Expense associado a user asset income existente. A procuedure irá criar um novo record na table FinancialuserItem e replicar (baseado no income existent), todos os records, associado o novo expense com as mesma datas.
STORED PROCEDURE TEST CASE(S):

CALL CreateUserAssetChildIncomeTax(1, 2,'New Child Tax IRPJ',3,60,950,'');

BACKEND VISUALIZATION:
  
select * from financialuseritem where financialuseritemname='New Child Tax IRPJ' -- mostrar o novo record (associado com o seu parent income)
select * from userfinancialforecast where financialuseritemid = (select financialuseritemid from financialuseritem where financialuseritemname='New Child Tax IRPJ') -- Mostrar todos os records na forecast


USER INTERFACE:

TBD
----------------------------------------------------*/
LANGUAGE plpgsql
AS $$

DECLARE 
    v_UserValidation INT;
    v_FinancialUserItemID INT;
    v_ParentRecurrency INT;
BEGIN
    -- Verifica se algum campo obrigatório está nulo ou vazio
    IF p_UserID IS NULL THEN
        p_Message := '{"status": "fail", "message": "Missing UserID"}';
        RETURN;
    END IF;

    IF p_ParentFinancialUserItemID IS NULL THEN
        p_Message := '{"status": "fail", "message": "Parent Income Item is mandatory"}';
        RETURN;
    END IF;

    -- Validação do UserAssetID
    SELECT userprofileid  
    INTO v_UserValidation
    FROM userasset
    WHERE userassetid = p_UserAssetID;
    

    IF p_UserID <> v_UserValidation THEN
        p_Message := '{"status": "fail", "message": "UserID provided is different from the asset owner"}';
        RETURN;
    END IF;

    -- Insere o novo Child Income Expense
    --Replicando mesma recorrencia do Parent Income
    SELECT recurrencyid  
    INTO v_ParentRecurrency
    FROM FinancialUserItemName
    WHERE FinancialUserItem = p_ParentFinancialUserItemID;
    INSERT INTO FinancialUserItem (
        FinancialUserItemName, EntityID, UserEntityID, RecurrencyID, FinancialUserEntityItemID, ParentFinancialUserItemID
    ) 
    VALUES (
        p_FinancialUserItemName, 13, p_UserAssetID, v_ParentRecurrency, p_FinancialUserEntityItemID, p_ParentFinancialUserItemID
    )
    RETURNING FinancialUserItemID INTO v_FinancialUserItemID;

    -- Criação dos valores associados ao forecast
    BEGIN
        -- Para cada forecast existente, crie um novo, substituindo os valores conforme necessário
        INSERT INTO userfinancialforecast (
            usercategoryid, financialuseritemid, userfinancialforecastbegindate, 
            userfinancialforecastenddate, userfinancialforecastamount, currencyid
        )
        SELECT 
            usercategoryid, 
            v_FinancialUserItemID,  -- Novo FinancialUserItemID
            userfinancialforecastbegindate, 
            userfinancialforecastenddate, 
            p_ExpenseAmount,  -- Valor da despesa
            currencyid
        FROM userfinancialforecast
        WHERE financialuseritemid = p_ParentFinancialUserItemID;

        -- Se não houver erro, envia mensagem de sucesso
        p_Message := '{"status": "success", "message": "New user expense associated to the parent income created successfully."}';
    EXCEPTION
        WHEN OTHERS THEN
            -- Em caso de erro, reverte as mudanças e retorna erro
            ROLLBACK;
            p_Message := '{"status": "fail", "message": "Could not create the Forecast, rolling back changes"}';
    END;

END;
$$;

-- DELETE User Asset Child Income Expense
CREATE OR REPLACE PROCEDURE DeleteUserAssetChildIncomeExpense(
    IN p_FinancialUserItemID INT,
    IN p_UserID INT,
    IN p_UserAssetID INT,
    OUT p_Message TEXT
)
/* ----------------------------------------------------------------------
STORED PROCEDURE NAME: DeleteUserAssetChildIncomeExpense
STORED PROCEDURE VERSION: 1.0
STORED PROCEDURE LAST UPDATED DATE: 30-Mar-2025
STORED PROCEDURE DESCRIPTION: 
     Deletar um expense child User Asset Income. Essa procedure é para um hard delete, excluindo todo historico tanto de forecasts e de actuals. Não existe opção de soft delete em childs Aqui todos os dados são deletados sem possibilidade de recuperação.
STORED PROCEDURE TEST CASE(S):

CALL CreateUserAssetChildIncomeTax(60,1,2,'')

BACKEND VISUALIZATION:
  
select * from financialuseritem where FinancialUserItemID=60

--Deve retorna "0" results em todos selects acima

USER INTERFACE:

TBD
----------------------------------------------------*/
LANGUAGE plpgsql
AS $$
DECLARE 
    v_UserValidation INT;
    v_EntityID INT;
BEGIN
    -- Verifica se o FinancialUserItem existe e obtém seu EntityID
    SELECT EntityID, UserEntityID INTO v_EntityID, v_UserValidation
    FROM FinancialUserItem
    WHERE FinancialUserItemID = p_FinancialUserItemID;
    
    -- Se o item não existir, retorna erro
    IF v_EntityID IS NULL THEN
        p_Message := '{"status": "fail", "message": "Child Income Expense not found"}';
        RETURN;
    END IF;
    
    -- Verifica se o EntityID é 13 (Expense)
    IF v_EntityID <> 13 THEN
        p_Message := '{"status": "fail", "message": "Not an expense type child"}';
        RETURN;
    END IF;
    
    -- Verifica se o UserAsset pertence ao UserID informado
    IF p_UserID <> v_UserValidation THEN
        p_Message := '{"status": "fail", "message": "UserID provided is different from the asset owner"}';
        RETURN;
    END IF;
	-- Deletar Relation records
    DELETE FROM UserForecastActualRelation
    WHERE UserFinancialForecastID IN (SELECT UserFinancialForecastID FROM userfinancialforecast where financialuseritemid=p_FinancialUserItemID)
    OR UserFinancialActualID IN (SELECT UserFinancialActualID FROM userfinancialactual where financialuseritemid=p_FinancialUserItemID);
        
    -- Deletar registros relacionados na tabela UserFinancialForecast
    DELETE FROM UserFinancialForecast
    WHERE FinancialUserItemID = p_FinancialUserItemID;
    
    -- Deletar o FinancialUserItem
    DELETE FROM FinancialUserItem 
    WHERE FinancialUserItemID = p_FinancialUserItemID;
    
    -- Retorna mensagem de sucesso
    p_Message := '{"status": "success", "message": "Child Asset Income Expense deleted successfully."}';
    
EXCEPTION 
    WHEN OTHERS THEN 
        p_Message := format(
            '{"status": "fail", "message": "An error occurred: %s"}', 
            SQLERRM
        );
END;
$$;





-- DELETE User Asset Child Income Tax

CREATE OR REPLACE PROCEDURE DeleteUserAssetChildIncomeTax(
    IN p_FinancialUserItemID INT,
    IN p_UserID INT,
    IN p_UserAssetID INT,
    OUT p_Message TEXT
)
/* ----------------------------------------------------------------------
STORED PROCEDURE NAME: DeleteUserAssetChildIncomeTax
STORED PROCEDURE VERSION: 1.0
STORED PROCEDURE LAST UPDATED DATE: 30-Mar-2025
STORED PROCEDURE DESCRIPTION: 
     Deletar um tax child User Asset Income. Essa procedure é para um hard delete, excluindo todo historico tanto de forecasts e de actuals. Não existe opção de soft delete em childs Aqui todos os dados são deletados sem possibilidade de recuperação.
STORED PROCEDURE TEST CASE(S):

CALL CreateUserAssetChildIncomeTax(60,1,2,'')

BACKEND VISUALIZATION:
  
select * from financialuseritem where FinancialUserItemID=60

--Deve retorna "0" results em todos selects acima

USER INTERFACE:

TBD
----------------------------------------------------*/
LANGUAGE plpgsql
AS $$
DECLARE 
    v_UserValidation INT;
    v_EntityID INT;
BEGIN
    -- Verifica se o FinancialUserItem existe e obtém seu EntityID
    SELECT EntityID, UserEntityID INTO v_EntityID, v_UserValidation
    FROM FinancialUserItem
    WHERE FinancialUserItemID = p_FinancialUserItemID;
    
    -- Se o item não existir, retorna erro
    IF v_EntityID IS NULL THEN
        p_Message := '{"status": "fail", "message": "Child Income Tax not found"}';
        RETURN;
    END IF;
    
    -- Verifica se o EntityID é 12 (Expense)
    IF v_EntityID <> 12 THEN
        p_Message := '{"status": "fail", "message": "Not an Tax type child"}';
        RETURN;
    END IF;
    
    -- Verifica se o UserAsset pertence ao UserID informado
    IF p_UserID <> v_UserValidation THEN
        p_Message := '{"status": "fail", "message": "UserID provided is different from the asset owner"}';
        RETURN;
    END IF;
	-- Deletar Relation records
    DELETE FROM UserForecastActualRelation
    WHERE UserFinancialForecastID IN (SELECT UserFinancialForecastID FROM userfinancialforecast where financialuseritemid=p_FinancialUserItemID)
    OR UserFinancialActualID IN (SELECT UserFinancialActualID FROM userfinancialactual where financialuseritemid=p_FinancialUserItemID);
        
    -- Deletar registros relacionados na tabela UserFinancialForecast
    DELETE FROM UserFinancialForecast
    WHERE FinancialUserItemID = p_FinancialUserItemID;
    
    -- Deletar o FinancialUserItem
    DELETE FROM FinancialUserItem 
    WHERE FinancialUserItemID = p_FinancialUserItemID;
    
    -- Retorna mensagem de sucesso
    p_Message := '{"status": "success", "message": "Child Asset Income Tax deleted successfully."}';
    
EXCEPTION 
    WHEN OTHERS THEN 
        p_Message := format(
            '{"status": "fail", "message": "An error occurred: %s"}', 
            SQLERRM
        );
END;
$$;

-- STORED PROCEDURE TO CREATE NEW PARENT EXPENSE
CREATE OR REPLACE PROCEDURE CreateUserParentExpense(
    IN p_UserID INT,
    IN p_FinancialUserItemName VARCHAR(255),
    IN p_RecurrencyID INT,
    IN p_FinancialUserEntityItemID INT,
    IN p_ParentExpenseAmount NUMERIC(15,2),
    IN p_BeginDate DATE,
    OUT p_Message TEXT
)

/* ----------------------------------------------------------------------
STORED PROCEDURE NAME: CreateUserParentExpense
STORED PROCEDURE VERSION: 1.0
STORED PROCEDURE LAST UPDATED DATE: 30-Mar-2025
STORED PROCEDURE DESCRIPTION: 
    Procedure para criar um novo parent expense para o usuario. Associado na tabela financialuseritem o novo record e criando (baseado na recurency escolhida) todos os records the forecast associados.
STORED PROCEDURE TEST CASE(S):

CALL CreateUserParentExpense (13,'Test Expense One Time-01',1,1,30,'05-03-2025','')  -- Para Criação de "One Time" expenses
CALL CreateUserParentExpense (13,'Test Expense Monthly-01',2,2,100,'05-03-2025','') -- Para Criação de "Monthly" expenses
CALL CreateUserParentExpense (13,'Test Expense Quarterly-01',3,4,1000,'05-03-2025','')-- Para Criação de "Quarterly" expenses
CALL CreateUserParentExpense (13,'Test Expense Yearly-01',4,5,2500,'05-03-2025','') -- Para Criação de "Time Yearly" expenses

BACKEND VISUALIZATION:
  
select fo.userfinancialforecastid,F.financialuseritemid,U.UserProfileId,F.userentityid,U.FirstName,U.LastName,F.financialuseritemname,E.EntityName,E.entitytype,FO.userfinancialforecastamount,FO.userfinancialforecastbegindate,FO.userfinancialforecastenddate from UserProfile U
Join financialuseritem F on F.userentityid=U.UserProfileID
right join userfinancialforecast FO on FO.financialuseritemid=F.financialuseritemid
join entity E ON E.entityid=F.entityid
where F.userentityid=13  -- Selecione o UserProfileID do usuario que foi criado os expense
and E.EntityID=6
-- and f.financialuseritemid=59 -- Se quiser granularidade no item que foi criado, use o financialuseritemid associado ao expense
order by 1,10

USER INTERFACE:
Navegue até Expense, e veja os records no sitema
----------------------------------------------------*/
LANGUAGE plpgsql
AS $$
DECLARE 
    v_NewFinancialUserItemid INT;
    v_CurrentDate DATE;
    v_NextDate DATE;
    v_Iterations INT;
    v_Increment INTERVAL;
    i INT;
BEGIN
    -- Validações de campos obrigatórios
    IF p_UserID IS NULL THEN 
        p_Message := '{"status": "fail", "message": "Missing UserID"}';
        RETURN;
    END IF;
    IF p_FinancialUserItemName IS NULL OR p_FinancialUserItemName = '' THEN 
        p_Message := '{"status": "fail", "message": "Missing FinancialUserItemName"}';
        RETURN;
    END IF;
    IF p_RecurrencyID IS NULL THEN 
        p_Message := '{"status": "fail", "message": "Missing RecurrencyID"}';
        RETURN;
    END IF;
    IF p_FinancialUserEntityItemID IS NULL THEN 
        p_Message := '{"status": "fail", "message": "Missing FinancialUserEntityItemID"}';
        RETURN;
    END IF;
    IF p_ParentExpenseAmount IS NULL OR p_ParentExpenseAmount <= 0 THEN 
        p_Message := '{"status": "fail", "message": "Invalid p_ParentExpenseAmount"}';
        RETURN;
    END IF;
    IF p_BeginDate IS NULL THEN 
        p_Message := '{"status": "fail", "message": "Missing BeginDate"}';
        RETURN;
    END IF;

    -- Inicia transação manualmente
    BEGIN
        -- Insere um novo Parent Expense
        INSERT INTO FinancialUserItem (
            FinancialUserItemName, EntityID, UserEntityID, RecurrencyID, 
            FinancialUserEntityItemID, ParentFinancialUserItemID
        ) VALUES (
            p_FinancialUserItemName, 6, p_UserID, p_RecurrencyID, 
            p_FinancialUserEntityItemID, NULL
        ) RETURNING FinancialUserItemID INTO v_NewFinancialUserItemid;

        -- Configura número de inserções e intervalo conforme a recorrência
        IF p_RecurrencyID = 1 THEN -- One time
            v_Iterations := 1;
            v_Increment := '1 day'::INTERVAL;
        ELSIF p_RecurrencyID = 2 THEN -- Monthly
            v_Iterations := 12;
            v_Increment := '1 month'::INTERVAL;
        ELSIF p_RecurrencyID = 3 THEN -- Quarterly
            v_Iterations := 4;
            v_Increment := '4 months'::INTERVAL;
        ELSIF p_RecurrencyID = 4 THEN -- Yearly
            v_Iterations := 1;
            v_Increment := '1 year'::INTERVAL;
        ELSE
            RAISE EXCEPTION 'Invalid RecurrencyID';
        END IF;

        -- Insere múltiplos registros conforme a recorrência
        v_CurrentDate := p_BeginDate;
        FOR i IN 1..v_Iterations LOOP
            -- Define a data do próximo registro para calcular a end date
            v_NextDate := v_CurrentDate + v_Increment;
            
            INSERT INTO UserFinancialForecast (
                usercategoryid, financialuseritemid, userfinancialforecastbegindate, 
                userfinancialforecastenddate, userfinancialforecastamount, currencyid
            ) VALUES (
                NULL, v_NewFinancialUserItemid, v_CurrentDate, 
                CASE 
                    WHEN p_RecurrencyID = 1 THEN v_CurrentDate + INTERVAL '1 day' - INTERVAL '1 day' -- One Time
                    WHEN p_RecurrencyID = 2 AND i < v_Iterations THEN v_NextDate - INTERVAL '1 day' -- Monthly
                    WHEN p_RecurrencyID = 2 AND i = v_Iterations THEN (v_CurrentDate + INTERVAL '1 month') - INTERVAL '1 day' -- Último mês
                    WHEN p_RecurrencyID = 3 AND i < v_Iterations THEN v_NextDate - INTERVAL '1 day' -- Quarterly
                    WHEN p_RecurrencyID = 3 AND i = v_Iterations THEN (v_CurrentDate + INTERVAL '4 months') - INTERVAL '1 day' -- Último trimestre
                    WHEN p_RecurrencyID = 4 THEN v_CurrentDate + INTERVAL '1 year' - INTERVAL '1 day' -- Yearly
                END,
                p_ParentExpenseAmount, 1
            );
            
            -- Atualiza data para a próxima recorrência
            v_CurrentDate := v_NextDate;
        END LOOP;

        -- Se tudo deu certo, define mensagem de sucesso
        p_Message := '{"status": "success", "message": "New user expense and forecast created successfully."}';
    EXCEPTION 
        WHEN OTHERS THEN
            -- Captura erro e define mensagem de falha
            p_Message := '{"status": "fail", "message": "Error creating forecast values: ' || SQLERRM || '"}';
    END;
END;
$$;

CREATE OR REPLACE PROCEDURE UpdateUserParentExpense(
    IN p_FinancialUserItemID INT,
    IN p_UserID INT,
    IN p_NewFinancialUserItemName VARCHAR(255),
    IN p_NewParentExpenseAmount NUMERIC(15,2),
    IN p_NewBeginDate DATE,
    IN p_IsActive BOOLEAN,
    OUT p_Message TEXT
)
/* ----------------------------------------------------------------------
STORED PROCEDURE NAME: UpdateUserParentExpense
STORED PROCEDURE VERSION: 1.0
STORED PROCEDURE LAST UPDATED DATE: 31-Mar-2025
STORED PROCEDURE DESCRIPTION: 
   Procedure altera valores de Expense baseados na data (utilizados para modificar um valor de um Expense a partir de uma data para frente no forecast)
   Tambe utilizada para inativar (Expense não é mais valido, porém mantem o historico) pela Flag IsActive
STORED PROCEDURE TEST CASE(S):

CALL UpdateUserParentExpense (66,13,'Updated Expense Name',3500,'2025-04-01',TRUE,'') -- Update Expense Value
CALL UpdateUserParentExpense (67,13,'Updated Expense Recurring',600.50,'2025-10-01',TRUE,'') -- Update Recurring value (Expensier)
CALL UpdateUserParentExpense (67,13,'Updated Expense Recurring',600.50,'2025-12-01',FALSE,'') -- Removing the Expense (Active=False) after a given date
BACKEND VISUALIZATION:
  
select fo.userfinancialforecastid,F.financialuseritemid,U.UserProfileId,F.userentityid,U.FirstName,U.LastName,F.financialuseritemname,E.EntityName,E.entitytype,FO.userfinancialforecastamount,FO.userfinancialforecastbegindate,FO.userfinancialforecastenddate from UserProfile U
Join financialuseritem F on F.userentityid=U.UserProfileID
right join userfinancialforecast FO on FO.financialuseritemid=F.financialuseritemid
join entity E ON E.entityid=F.entityid
where F.userentityid=13  -- Selecione o UserProfileID do usuario que foi criado os Expenses
-- and f.financialuseritemid=59 -- Se quiser granularidade no item que foi criado, use o financialuseritemid associado ao Expense
order by 1,10

USER INTERFACE:
Navegue até Expense, e veja os records no sitema
----------------------------------------------------*/
LANGUAGE plpgsql
AS $$
DECLARE 
    v_FinancialItemExist INT;
    v_AssociatedActual INT;
BEGIN

-- Validações de campos obrigatórios
    IF p_FinancialUserItemID IS NULL THEN 
        p_Message := '{"status": "fail", "message": "Missing FinancialUserItemID"}';
        RETURN;
    END IF;
    IF p_UserID IS NULL THEN 
        p_Message := '{"status": "fail", "message": "Missing UserID"}';
        RETURN;
    END IF;
    IF p_NewFinancialUserItemName IS NULL OR p_NewFinancialUserItemName = '' THEN 
        p_Message := '{"status": "fail", "message": "Missing the new FinancialUserItemName"}';
        RETURN;
    END IF;
    IF p_NewParentExpenseAmount IS NULL THEN 
        p_Message := '{"status": "fail", "message": "Missing the new ParentExpenseAmount"}';
        RETURN;
    END IF;
    IF p_NewBeginDate IS NULL THEN 
       p_Message := '{"status": "fail", "message": "Missing new BeginDate"}';
       RETURN;
    END IF;
    -- Verifica se o FinancialUserItem existe
    SELECT financialuseritemid 
    INTO v_FinancialItemExist
    FROM UserFinancialForecast
    WHERE financialuseritemid = p_FinancialUserItemID
    LIMIT 1;     
    IF v_FinancialItemExist IS NULL THEN
        p_Message := '{"status": "fail", "message": "UserParentExpense not found"}';
        RETURN;
    END IF;
    -- Confere se o user realmente é o dono do financial user item
    IF p_userID <> (SELECT userentityid FROM FinancialUserItem WHERE FinancialUserItemID = p_FinancialUserItemID) THEN
       p_Message := '{"status": "fail", "message": "User do not match FinancialuserItem"}';
       RETURN;
    END IF;

    IF p_isActive = TRUE THEN
        -- Atualiza o nome do FinancialUserItem para todos os records
        IF p_NewFinancialUserItemName IS NOT NULL AND p_NewFinancialUserItemName <> '' THEN
            UPDATE FinancialUserItem
            SET FinancialUserItemName = p_NewFinancialUserItemName
            WHERE FinancialUserItemID = p_FinancialUserItemID;
        END IF;

        -- Atualiza apenas os valores das previsões futuras
        IF p_NewBeginDate IS NOT NULL AND p_NewParentExpenseAmount IS NOT NULL THEN
            UPDATE UserFinancialForecast
            SET userfinancialforecastamount = p_NewParentExpenseAmount
            WHERE financialuseritemid = p_FinancialUserItemID
            AND userfinancialforecastbegindate >= p_NewBeginDate;
        END IF;

        p_Message := '{"status": "success", "message": "UserParentExpense updated successfully."}';
    END IF;

    IF p_isActive = FALSE THEN
    -- Veja se existe algum Forecast com actuals associado para esse item com dados no futuro
    
        SELECT FROM userforecastactualrelation fr INTO v_AssociatedActual
        WHERE fr.userfinancialforecastid in (SELECT UserFinancialForecastid FROM UserFinancialForecast f
        WHERE f.FinancialUserItemID=p_FinancialUserItemID and f.userfinancialforecastbegindate > p_NewBeginDate)
        LIMIT 1;
        IF v_AssociatedActual IS NULL THEN    
        -- Inativando o FinancialUserItem
        UPDATE FinancialUserItem
        SET IsActive = FALSE
        WHERE financialuseritemid = p_FinancialUserItemID;
        
        -- Deletar todos os forecasts que ainda existem no futuro da data fornecida
        DELETE FROM UserFinancialForecast
        WHERE financialuseritemid = p_FinancialUserItemID
        AND userfinancialforecastbegindate > p_NewBeginDate;
        
        p_Message := '{"status": "success", "message": "Forecasting for this UserParentExpense is inactive"}';
        ELSE 
        -- Informar que precisa ou mudar a data ou remover o actual associado
         p_Message := '{"status": "fail", "message": "There is an actual record associated with a forecast in a future data provided, delete the actual or adjust the date, inactivation aborted"}';
        END IF;
    END IF;

END;
$$;

-- DELETE PARENT INCOME
CREATE OR REPLACE PROCEDURE DeleteUserParentExpense(
    IN p_FinancialUserItemID INT,
    IN p_UserID INT,
    OUT p_Message TEXT
)
/* ----------------------------------------------------------------------
STORED PROCEDURE NAME: DeleteUserParentExpense
STORED PROCEDURE VERSION: 1.0
STORED PROCEDURE LAST UPDATED DATE: 31-Mar-2025
STORED PROCEDURE DESCRIPTION: 
   Deletar o Parent Expense. Essa procedure é para um hard delete, excluindo todo historico tanto de forecasts e de actuals (soft deletes são executados pelo UpdateUserParentExpense, pela flag IsActive). Aqui todos os dados são deletados sem possibilidade de recuperação.
   Todos os items relacionado a user Parent serão deletados.
STORED PROCEDURE TEST CASE(S):

call DeleteUserParentExpense (66,13,'')-- Informando FinancialUserItemID e UserID

BACKEND VISUALIZATION:
  
select * from userfinancialactual where FinancialUserItemID=60 -- Selectione o FinancialUserItemID no select
select * from u userfinancialforecast where FinancialUserItemID=60 -- Selectione o FinancialUserItemID no select
select * from financialuseritem where FinancialUserItemID=60 -- Selectione o FinancialUserItemID no select
Deve retorna "0" results em todos selects acima-- Selectione o FinancialUserItemID no select

USER INTERFACE:

Não devem estar visiveis na UI após o delete.
----------------------------------------------------*/
LANGUAGE plpgsql
AS $$
DECLARE 
    v_UserValidation INT;
    v_ChildCheck INT;
BEGIN
    -- Verifica se o FinancialUserItem existe
    IF NOT EXISTS (
        SELECT 1 FROM FinancialUserItem 
        WHERE FinancialUserItemID = p_FinancialUserItemID 
        AND UserEntityID = p_UserID
    ) THEN
        p_Message := '{"status": "fail", "message": "UserParentExpense not found"}';
        RETURN;
    END IF;
     
        -- Deletar o User Parent Expense da relations
        DELETE FROM UserForecastActualRelation
        WHERE UserFinancialForecastID IN (SELECT UserFinancialForecastID FROM userfinancialforecast where financialuseritemid=p_FinancialUserItemID)
        OR UserFinancialActualID IN (SELECT UserFinancialActualID FROM userfinancialactual where financialuseritemid=p_FinancialUserItemID);
        
        -- Deletar o User Parent Expense da Forecast
        DELETE FROM UserFinancialForecast
        WHERE financialuseritemid =p_FinancialUserItemID;

        
        -- Deletar o User Parent Expense da Actual
        DELETE FROM FinancialUserItem 
        WHERE FinancialUserItemID = p_FinancialUserItemID;

        -- Retorna mensagem de sucesso
        p_Message := '{"status": "success", "message": "User Parent Expense and all references deleted successfully."}';

        EXCEPTION 
        WHEN OTHERS THEN 
            -- Se algum erro ocorrer, retorna uma mensagem de erro e encerra a procedure
            p_Message := format(
                '{"status": "fail", "message": "An error occurred: %s"}', 
                SQLERRM
            );
END;
$$;